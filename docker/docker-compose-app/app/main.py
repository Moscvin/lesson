from pymongo import MongoClient # type: ignore
from pprint import pprint

MONGO_URL = 'mongodb://mongo:27017'
client = MongoClient(MONGO_URL)
db = client.admin
dbs_list = db.command("listDatabases")
pprint(dbs_list)

print('Hello World!')