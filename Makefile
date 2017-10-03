BIN = cleopatchra

$(BIN):
	$(CC) -std=c99 -lpq $(@).c -o $(@)
	./$(BIN)

clean:
	rm -f $(BIN)

.PHONY: $(BIN)