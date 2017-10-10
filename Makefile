BIN = cleopatchra
CC = g++

$(BIN):
	$(CC) -std=c++11 main.cpp -o $(@)
	./$(BIN)

clean:
	rm -f $(BIN)

.PHONY: $(BIN)