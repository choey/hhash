wordnet_dir=./db
sample_dir=./db/sample
datasets_go_template=./datasets.go.template

generate-datasets:
	@echo "[+] generating smaller datasets from WordNet database"
	@./script/generate-small-dataset.sh -d $(wordnet_dir) -t $(sample_dir)
	@echo "[+] generating datasets.go"
	@./script/generate-default-datasets.sh -d $(sample_dir) -t $(datasets_go_template)
	@echo "[!] done generating datasets"