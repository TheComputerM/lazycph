.PHONY: dev
dev:
	@uv run textual run --dev src/lazycph/__main__.py

.PHONY: serve
serve:
	@uv run textual serve src/lazycph/__main__.py

.PHONY: inspect
inspect:
	@uv run textual console -x SYSTEM -x DEBUG -x EVENT

.PHONY: test
test:
	@uv run pytest tests/ -v

.PHONY: build
build:
	@uv run pyinstaller lazycph.spec
