
# Operations API Examples

## Plant a Crop

This endpoint allows a user to plant a crop on a leased plot.

**Request:**

```bash
curl -X POST -H "Content-Type: application/json" -d '{"crop_id": "<CROP_ID>"}' http://localhost:8080/api/v1/plots/<PLOT_ID>/plantings
```

**Response:**

```json
{
  "id": "<PLANTING_ID>",
  "plot_id": "<PLOT_ID>",
  "crop_id": "<CROP_ID>",
  "lease_id": "<LEASE_ID>",
  "planted_at": "<TIMESTAMP>",
  "status": "growing",
  "created_at": "<TIMESTAMP>",
  "updated_at": "<TIMESTAMP>"
}
```

## Remove a Crop

This endpoint allows a user to remove a planted crop from a plot.

**Request:**

```bash
curl -X DELETE http://localhost:8080/api/v1/plots/<PLOT_ID>/plantings/<PLANTING_ID>
```

**Response:**

```
204 No Content
```
