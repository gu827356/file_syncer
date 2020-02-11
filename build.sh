outdir=bin
rm -rf $outdir

go build -o $outdir/file_syncer_client ./cli/client
go build -o $outdir/file_syncer_server ./cli/server