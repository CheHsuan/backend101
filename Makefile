.PHONY: static clean

static:
	markdown-tree --input-dir ../backend101 --output-file tree.md --with-link --base-url https://github.com/CheHsuan/backend101/tree/main/
	npx markmap-cli --no-open -o static/index.html tree.md

clean:
	rm tree.md
