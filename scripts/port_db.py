import pymongo
import json
from tqdm import tqdm
from time import sleep
import requests
import datetime

index = {"index": {"_index": "chatlogs_foo", "_type": "_doc"}}


# Format for the .bson documents:
# {
#   "Messages": {
#       "Player": {"UUID": "the-uuid-bla-bla"},
#       "Time": "the-timestamp-as-unix-long"},
#       "Message": "foo"
#   }
# }
def _process_documents(docs):
    data = ""

    for doc in docs:
        msgs = doc["Messages"]
        for msg in msgs:
            uuid = str(msg["Player"])
            raw_msg = msg["Message"]
            time = int(msg["Time"])

            d = datetime.datetime.fromtimestamp(time / 1e3)
            data += json.dumps(index) + "\n"
            data += json.dumps({"user": uuid, "message": raw_msg, "time": d.isoformat("T") + "Z"}) + "\n"

    return requests.post("http://localhost:9200/_bulk/", data=data.encode("utf-8"),
                         headers={"Content-Type": "application/x-ndjson"})


# the url and credentials of the mongo database
url = "127.0.0.1/test"
full_url = "mongodb://user:password@" + url

print("Connect to mongodb @" + url + " ...")
client = pymongo.MongoClient(full_url)
print("Connected to mongodb.")

database = client['test']
collection = database['chatlogs']

doc_size = collection.estimated_document_count()
print("Fetched collection " + collection.name + " with {} documents".format(doc_size))

# the maximum size to be fetched from
# the mongodb
max_size = 250000

# size of documents inside a bulk update
part_size = 1000

# create the progress bar with 100%=max_size
bar = tqdm(total=max_size)

buffer = []
cursor = collection.find()
count = 0

for i in range(max_size):
    # if the buffer has the max size of a bulk write
    # process the documents before continuing
    if len(buffer) == part_size:
        count = 0
        result = _process_documents(buffer)

        if not result.ok:
            print("Unexpected result: " + str(result))
            break
        buffer.clear()
        sleep(0.1)
    # update the progress bar
    bar.update(1)

    # add the document to the buffer
    buffer.append(cursor.next())

sleep(1)
print("Done.")
