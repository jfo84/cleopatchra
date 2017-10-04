#include <iostream>
#include <libpq-fe.h>
#include "./lib/picojson.h"

using namespace std;

int main()
{
  PGconn *conn;

  conn = connect("127.0.0.1", "5432", "cleopatchra");

  if (NULL == conn) { return 1; }

  // TODO: Figure out how to query

  // query
  PGresult *res = PQexec(conn, query);

  return 0;
}

PGconn *connect(const char *hostaddr, const char *port, const char *dbname)
{
  const char *keys[] = {
    "hostaddr", "port", "dbname", NULL
  };

  const char *values[] = {
    hostaddr, port, dbname, NULL
  };

  PGconn *conn = PQconnectdbParams(keys, values, 0);
  if (NULL == conn) { return NULL; }

  ConnStatusType status = PQstatus(conn);

  if (CONNECTION_BAD == status) {
    fprintf(stderr, "error: %s\n", PQerrorMessage(conn));
    PQfinish(conn);
    return NULL;
  }

  return conn;
}