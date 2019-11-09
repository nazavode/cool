TMJAR:=$(if $(TEXTMAPPER),$(TEXTMAPPER),textmapper.jar)

all: gen

gen: cool.tm
	@java -jar ${TMJAR} $<
	@go fmt ./... > /dev/null
	@go build ./...

clean:
	$(RM) -v cool listener.go lexer.go lexer_tables.go parser.go parser_tables.go token.go
	$(RM) -rf -v ast/
	$(RM) -rf -v selector/

.PHONY: all gen clean