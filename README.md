# Video Meta Service
- REST service for updating meta
- Kafka Consumer for listening to video delete events (non-functional)

## Running the service:
```
make run
```

## Building docker image
```
make docker-build
```

## Implementation
Uses `AWS` implementation of tagging to allow user-centric metadata to be applied
to videos.

## API
_for more detailed and opertational calls see `example.sh`_
- `POST http://localhost:8080/<video_id> <payload>` -> CREATE meta data for a video
**payload**:
```json
[
    {
      "key": "<some_key>",
      "value": "<some_value>"
    },
    {
      "key": "<some_key>",
      "value": "<some_value>"
    }
]
```
- `GET http://localhost:8080/<video_id>` -> READ all meta for a video
- `PUT http://localhost:8080/<meta_id> <payload>` -> UPDATE meta data for a video
**payload**:
```json
{
  "key": "<some_key>",
  "value": "<some_value>"
}
```
- `DELETE http://localhost:8080/<meta_id>` -> DELETE meta data for a video
- Deleting of all information for a video would be done via listening for 
an event by a kafka consumer see `./consumer/consumer.go`

## Authorisation
All requests require the following headers:
```
X-API-KEY: SUPER_SECRET
Content-Type: application/json
```

## Data model
Is intentionally extremely flexible
