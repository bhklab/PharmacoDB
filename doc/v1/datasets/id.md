# Dataset

```
GET /datasets/{id}
```

## Description

This method returns a single dataset.

## Summary

| Name | Value |
| --- | --- |
| **Request Protocol** | GET |
| **Requires API Key** | No |
| **Cache Time** | 0 seconds |

## Notes

- A `404` error is returned if an item is not found.

## Sources

- http://pharmacodb.pmgenomics.ca/datasets

## Parameters

```
GET /datasets/{id}
```

| Parameter | Type | Value | Required | Default | Description |
| --- | --- | --- | --- | --- | --- |
| **indent** | input | *boolean* | no | false | Add indentation to response |
| **type** | input | - | no | dataset_id | Define whether `id = dataset_id` or `id = dataset_name` |

## Output Formats

- JSON

## Examples

```
GET /datasets/{id}
```

- https://api.pharmacodb.com/v1/datasets/1
- https://api.pharmacodb.com/v1/datasets/ccle?type=name
- https://api.pharmacodb.com/v1/datasets/1?indent=true

## Output

**JSON**, using the dataset **CCLE**.

```
{
    "id": 1,
    "name": "CCLE"
}
```
