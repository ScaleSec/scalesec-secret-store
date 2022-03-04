#!/bin/bash
#
# This shell sciprt is here instead of a make file to make it easy to build and install
# the plugin in Vault.
#
# USEAGE:
# make-scalesec-secret-store-plugin.sh <option: build; install; all>

##
## Funtion build
##
build () {
    echo "***************************************************"
    echo " BUILD"
    echo "***************************************************"

    set -e
    echo "go version *********"
    go version

    # Set the build information for the -ldflag
    BUILD_VERSION=$(date +%Y%m%d%H%M)
    BUILD_DATE=$(date +%Y-%m-%d)
    BUILT_INFORMATION="ScaleSec_Secret_Store_Plugin"

    echo "BUILT_INFORMATION=$BUILT_INFORMATION"
    echo "BUILD_DATE=$BUILD_DATE"
    echo "BUILD_VERSION=$BUILD_VERSION"

    # install gox
    # https://github.com/mitchellh/gox
    # we need to turn off go module so the gox will install  save the current setting
    # so we dont messup the current go environment and then restore it.
    if [[ ! -f ~/go/bin/gox ]]; then
        ORG_GO111MODULE=$GO111MODULE
        export GO111MODULE="off"   
        go install github.com/mitchellh/gox
        export GO111MODULE=$ORG_GO111MODULE
    fi

    # Make sure the current user go/bin is on the path.  This is where the gox is installed
    PATH_ORG=$PATH
    export PATH=$PATH:~/go/bin

    # clean up old files
    rm -rf bin
    rm -rf pkg
    mkdir -p bin/

    if [[ $DEBUG_FLAG == "TRUE" ]]; then
        echo "BUILDING for DEBUG  -gcflags \"all=-N -l\""
        gox \
            -verbose \
            -rebuild \
            -ldflags="-X main.pluginBuildDate=$BUILD_DATE -X main.pluginBuildVersion=$BUILD_VERSION -X main.pluginBuildInfo=$BUILT_INFORMATION" \
            -osarch="linux/amd64 darwin/amd64" \
            -output "bin/{{.OS}}_{{.Arch}}/scalesecSecretStorePlugin" \
            -gcflags "all=-N -l" \
            ./plugin/.
    else
        gox \
            -verbose \
            -rebuild \
            -ldflags="-X main.pluginBuildDate=$BUILD_DATE -X main.pluginBuildVersion=$BUILD_VERSION -X main.pluginBuildInfo=$BUILT_INFORMATION" \
            -osarch="linux/amd64 darwin/amd64" \
            -output "bin/{{.OS}}_{{.Arch}}/scalesecSecretStorePlugin" \
            ./plugin/.    
    fi
    echo ""
    echo "Results:"
    ls -hl ./bin/*

    #reset PATH to the way it was before we added the go/bin
    export PATH=$PATH_ORG
}

##
## deploy plugin to vault
##
deploy () {
    echo "***************************************************"
    echo " DEPLOY"
    echo "***************************************************"

    export VAULT_ROOT=~/vault 

    # Vault was installed and configured using the
    # ScaleSec Vault Assistant Project
    # GitHub: xxxxxxx

    # ScaleSec Vault Assistand Project makes it easy to run vault in a non-development mode.
    # Running in a non-development mode allows:
    # 1: Seal and Un-Seal Vault
    # 2: Persistant Store of Secrets
    # 3: Sgining of custom plugins
    
    #
    # make sure vault is configured for custom plugins
    #
    
    # Location of our custom plugins confgured in vault config.hcl
    VAULT_PLUGIN_DIR=~/vault/custom_plugin

    # Location of our config.hcl
    VAULT_CONFIG_HCL=~/vault/config.hcl

    if [[ ! -d $VAULT_PLUGIN_DIR ]]; then
        echo "Custom plugin directory not found $VAULT_PLUGIN_DIR ...  Consider using ScaleSec Vault Assistant"
        exit 9
    fi

    if [[ ! -f $VAULT_CONFIG_HCL ]]; then
        echo "Vault config.hcl not found at: $VAULT_CONFIG_HCL ... Consider using ScaleSec Vault Assistant"
        exit 9
    fi

    # make sure vault is configured to use the custom plugin
    if grep -q "$VAULT_PLUGIN_DIR" $VAULT_CONFIG_HCL; then
        echo "Vault HCL configured with custom plugin directory"
    else
        echo "Vault config.hcl missing: plugin_directory =\"$VAULT_PLUGIN_DIR\" ... Consider using ScaleSec Vault Assistant"
        exit 9
    fi

    #
    # Determin which OS plugin to install
    #

    if [[ "Darwin" == $(uname) ]] || [[ "darwin" == $(uname) ]]; then
        PLUGIN_TO_INSTALL=$PWD"/bin/darwin_amd64/scalesecSecretStorePlugin"
    else
        PLUGIN_TO_INSTALL=$PWD"/bin/linux_amd64/scalesecSecretStorePlugin"
    fi

    # Makesure vault is running
    if [[ `ps -ef | grep "vault/config.hcl" | grep -v grep` == "" ]]; then
        echo "Please start Vault and make sure it is unsealed ... Consider using ScaleSec Vault Assistant"
        exit 9
    fi
    
    # Vault is running make sure it is unsealed
    UNSEAL_STATUS=$(vault status | awk '/^Sealed/' | tr -s ' ' | cut -d ' ' -f2)
    if [[ $UNSEAL_STATUS == "false" ]]; then
       echo "Vault is Unsealed"
    else
        echo "Vault needs to be unsealed ... Consider using ScaleSec Vault Assistant"
        exit 9
    fi

    #
    vault_login

    # installing a plugin
    # https://www.vaultproject.io/docs/internals/plugins#plugin-catalog

    # Copy Pluign to custom_plugin directory
    cp $PLUGIN_TO_INSTALL $VAULT_PLUGIN_DIR

    # Enable mlock
    # sudo setcap cap_ipc_lock=+ep $VAULT_PLUGIN_DIR/scalesecSecretStorePlugin

    # Create the SHASUM for the plugin
    echo "Generate the shasum 256 for the plugin we are installing"
    SHASUM=$(shasum -a 256 $PLUGIN_TO_INSTALL | cut -d " "  -f1)
    echo "SHASUM=$SHASUM for $PLUGIN_TO_INSTALL"

    # unregister the plugin so the new verison will get picked up
    unregister

    # Register the plugin with Vault
    # new
    vault plugin register -sha256=$SHASUM secret scalesecSecretStorePlugin

    # old
    # vault write sys/plugins/catalog/secret/scalesecSecretStorePlugin sha256="$SHASUM" command=scalesecSecretStorePlugin

    # enable
    # vault secrets enable scalesecSecretStorePlugin
    vault secrets enable -path=scalesecsecrets scalesecSecretStorePlugin

    # Vault display plugin information
    vault plugin info secret/scalesecSecretStorePlugin

    # List secret eng
    vault secrets list

}

unregister () {
    echo "UNREGISTER"
    # disable and deregister so the new version is picked up
    PLUGIN_LIST=`vault secrets list | grep 'scalesecSecretStorePlugin'`
    if [[ $PLUGIN_LIST != "" ]]; then
        # vault secrets disable secrets at the path where we enabled the secrets
        vault secrets disable scalesecsecrets/
        # deregister the secret plugin scalesecSecretStorePlugin
        vault plugin deregister secret scalesecSecretStorePlugin
    fi

    # purge out old plugin processes
    PLUGIN_PROCCESS=`pgrep scalesecSecretStorePlugin`
    for PLUGIN_PROCESS_ID in $PLUGIN_PROCCESS
    do
        kill -9 $PLUGIN_PROCESS_ID
    done

}

vault_login () {
    # Get the root token
    VAULT_ROOT_TOKEN=""
    if [[ $VAULT_ROOT_TOKEN == "" ]]; then
        # if the token is not set manually then get we assume your using ScaleSec Vault Assistant
        # and will get the root token based on where the aisstant stortes it.
        if [[ ! -e ~/vault/local-root-token ]]; then 
            echo "Doesnt look like we were able to locate the root token.  Consider using ScaleSec Vault Assistant"
            exit 9
        else
            VAULT_ROOT_TOKEN=`cat ~/vault/local-root-token`
        fi
    fi
    # login to vault with the token
    vault login $VAULT_ROOT_TOKEN
}
debug () {
    export DEBUG_FLAG="TRUE"
}

all () {
    build
    deploy
}

test_write () {
    vault_login
    vault write scalesecsecrets/test secret_key="secret_value"
}

test_read () {
    vault_login
    vault read scalesecsecrets/test secret_key
}

test_delete () {
    vault_login
    vault delete scalesecsecrets/test secret_key
}

test () {
    test_write
    test_read
    test_delete
    
}

if [[ -z "$1" ]]; then
    all
    exit
fi

export OPTIONS=$@
for OPTION in $OPTIONS; do
    echo "OPTION=$OPTION"
   # if [[ $OPTION == "build" ]]; then
   #     build
   # fi
   #
   # if [[ $OPTION == "deploy" ]]; then
   #     deploy
   #43 fi

    case "$OPTION" in
        "debug") debug
        ;;
        "build") build
        ;;
        "deploy") deploy
        ;;
        "test") test
        ;;
        "test_write") test_write
        ;;
        "test_read") test_read
        ;; 
        "test_delete") test_delete
        ;;
        "test_delete") test_list
        ;;

        *) all
        ;;
    esac

done

