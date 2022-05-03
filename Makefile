.PHONY: static clean

static: clean
	markdown-tree --input-dir ../backend101 --output-file tree.md
	npx markmap-cli --no-open -o static/index.html tree.md

clean:
	rm tree.md
