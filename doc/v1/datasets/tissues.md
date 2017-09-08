# Tissues Tested In A Dataset

```
GET /datasets/{id}/tissues
```

## Description

This method returns a list of unique tissues that have been tested in a dataset of interest.

## Summary

| Name | Value |
| --- | --- |
| **Request Protocol** | GET |
| **Requires API Key** | No |
| **Cache Time** | 0 seconds |

## Notes

- Meta information is included in response headers by default. Use `include` parameter to add info to response body.
- A `404` error is returned if an item is not found.

## Sources

- http://pharmacodb.pmgenomics.ca/datasets

## Parameters

```
GET /datasets/{id}/tissues
```

| Parameter | Type | Value | Required | Default | Description |
| --- | --- | --- | --- | --- | --- |
| **page** | filter | *integer* | no | 1 | The page number for output |
| **per_page** | filter | *integer* | no | 30 | Number of items returned per page |
| **include** | input | metadata | no | - | Include meta info (eg. pagination) in body instead of headers |
| **indent** | input | *boolean* | no | false | Add indentation to response |
| **type** | input | - | no | dataset_id | Define whether `id = dataset_id` or `id = dataset_name` |

## Output Formats

- JSON

## Examples

```
GET /datasets/{id}/tissues
```

- https://api.pharmacodb.com/v1/datasets/1/tissues
- https://api.pharmacodb.com/v1/datasets/1/tissues?page=2&per_page=10
- https://api.pharmacodb.com/v1/datasets/ccle/tissues?type=name

## Output

**JSON**, using the dataset **CCLE**, and meta info included in body.

```
{
    "data": [
        {
            "id": 2,
            "name": "autonomic_ganglia"
        },
        {
            "id": 3,
            "name": "biliary_tract"
        },
        {
            "id": 5,
            "name": "bone"
        },
        {
            "id": 7,
            "name": "breast"
        },
        {
            "id": 9,
            "name": "central_nervous_system"
        },
        {
            "id": 11,
            "name": "endometrium"
        },
        {
            "id": 13,
            "name": "haematopoietic_and_lymphoid_tissue"
        },
        {
            "id": 15,
            "name": "kidney"
        },
        {
            "id": 16,
            "name": "large_intestine"
        },
        {
            "id": 17,
            "name": "liver"
        },
        {
            "id": 18,
            "name": "lung"
        },
        {
            "id": 21,
            "name": "oesophagus"
        },
        {
            "id": 24,
            "name": "ovary"
        },
        {
            "id": 25,
            "name": "pancreas"
        },
        {
            "id": 27,
            "name": "pleura"
        },
        {
            "id": 28,
            "name": "prostate"
        },
        {
            "id": 31,
            "name": "skin"
        },
        {
            "id": 33,
            "name": "soft_tissue"
        },
        {
            "id": 34,
            "name": "stomach"
        },
        {
            "id": 36,
            "name": "thyroid"
        },
        {
            "id": 37,
            "name": "upper_aerodigestive_tract"
        },
        {
            "id": 38,
            "name": "urinary_tract"
        }
    ],
    "metadata": {
        "last_page": 1,
        "page": 1,
        "per_page": 30,
        "total": 22
    }
}
```
