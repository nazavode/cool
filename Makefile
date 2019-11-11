GENSOURCES =	\
	token.go	\
	lexer.go	\
	lexer_tables.go 

all: gen

gen: $(GENSOURCES)

$(GENSOURCES): cool.tm
	@textmapper generate
	@go fmt ./... > /dev/null

clean:
	$(RM) $(GENSOURCES)
