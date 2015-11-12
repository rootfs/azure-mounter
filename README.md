# azure-mounter
create (if not exist yet) and mount share from Azure File Service

# Build and run
```console
go build -o azure-mounter main.go
./azure-mounter --account-name your_account_name --account-key your_key -mountpoint=/path_to_mountpoint --share=your_share_name_to_be_created
```