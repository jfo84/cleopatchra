BIN = cleopatchra
CC = g++

$(BIN):
	$(CC) -std=c++11 -I/usr/include/websocketpp/ main.cpp -o $(@)
	./$(BIN)

clean:
	rm -f $(BIN)

.PHONY: $(BIN)