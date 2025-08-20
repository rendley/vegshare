# Operations API Examples

## Plant a Crop

This endpoint allows a user to plant a crop on a leased plot.

**Request:**

```bash
ACCESS_TOKEN="your_access_token"
PLOT_ID="51175ae1-a6ae-45e2-9423-cce34fffcd63"
CROP_ID="62d71460-4689-4e3d-8e17-101ade9ab271"
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ACCESS_TOKEN" -d '{"crop_id": "$CROP_ID"}' http://localhost:8080/api/v1/operations/plots/$PLOT_ID/plantings
```

**Response:**

```json
{
  "id": "c09ffe51-fe12-4af1-bd32-0cc498399541",
  "plot_id": "51175ae1-a6ae-45e2-9423-cce34fffcd63",
  "crop_id": "62d71460-4689-4e3d-8e17-101ade9ab271",
  "lease_id": "4bbd645c-675c-49db-ab0e-618f4054ba77",
  "planted_at": "2025-08-20T20:41:26.137469084Z",
  "status": "growing",
  "created_at": "2025-08-20T20:41:26.137469084Z",
  "updated_at": "2025-08-20T20:41:26.137469084Z"
}
```

## Get Planted Crops on a Plot

This endpoint allows a user to get all the planted crops on a specific plot.

**Request:**

```bash
ACCESS_TOKEN="your_access_token"
PLOT_ID="51175ae1-a6ae-45e2-9423-cce34fffcd63"
curl -X GET -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:8080/api/v1/operations/plots/$PLOT_ID/plantings
```

**Response:**

```json
[
  {
    "id": "c09ffe51-fe12-4af1-bd32-0cc498399541",
    "plot_id": "51175ae1-a6ae-45e2-9423-cce34fffcd63",
    "crop_id": "62d71460-4689-4e3d-8e17-101ade9ab271",
    "lease_id": "4bbd645c-675c-49db-ab0e-618f4054ba77",
    "planted_at": "2025-08-20T20:41:26.137469Z",
    "status": "growing",
    "created_at": "2025-08-20T20:41:26.137469Z",
    "updated_at": "2025-08-20T20:41:26.137469Z"
  }
]
```

## Remove a Crop

This endpoint allows a user to remove a planted crop from a plot.

**Request:**

```bash
ACCESS_TOKEN="your_access_token"
PLOT_ID="51175ae1-a6ae-45e2-9423-cce34fffcd63"
PLANTING_ID="c09ffe51-fe12-4af1-bd32-0cc498399541"
curl -X DELETE -H "Authorization: Bearer $ACCESS_TOKEN" http://localhost:8080/api/v1/operations/plots/$PLOT_ID/plantings/$PLANTING_ID
```

**Response:**

```
204 No Content
```

## Perform an Action on a Plot

This endpoint allows a user to perform an action on a plot, such as watering.

**Request:**

```bash
ACCESS_TOKEN="your_access_token"
PLOT_ID="51175ae1-a6ae-45e2-9423-cce34fffcd63"
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $ACCESS_TOKEN" -d '{"action": "water"}' http://localhost:8080/api/v1/operations/plots/$PLOT_ID/actions
```

**Response:**

```json
{
  "message": "Action performed successfully"
}
```