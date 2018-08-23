# Compounds Tested In A Dataset

```
GET /datasets/{id}/compounds
```

## Description

This method returns a list of unique compounds that have been tested in a dataset of interest.

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

- http://pharmacodb.ca/datasets

## Parameters

```
GET /datasets/{id}/compounds
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
GET /datasets/{id}/compounds
```

- http://api.pharmacodb.ca/v1/datasets/1/compounds
- http://api.pharmacodb.ca/v1/datasets/1/compounds?page=2&per_page=10
- http://api.pharmacodb.ca/v1/datasets/ccle/compounds?type=name

## Output

**JSON**, using the dataset **CCLE**, and meta info included in body.

```
{
    "data": [
        {
            "id": 3,
            "name": "17-AAG"
        },
        {
            "id": 21,
            "name": "AEW541"
        },
        {
            "id": 53,
            "name": "AZD0530"
        },
        {
            "id": 56,
            "name": "AZD6244"
        },
        {
            "id": 248,
            "name": "Crizotinib"
        },
        {
            "id": 287,
            "name": "Erlotinib"
        },
        {
            "id": 367,
            "name": "Irinotecan"
        },
        {
            "id": 413,
            "name": "L-685458"
        },
        {
            "id": 415,
            "name": "lapatinib"
        },
        {
            "id": 417,
            "name": "LBW242"
        },
        {
            "id": 491,
            "name": "Nilotinib"
        },
        {
            "id": 507,
            "name": "Nutlin-3"
        },
        {
            "id": 526,
            "name": "paclitaxel"
        },
        {
            "id": 529,
            "name": "Panobinostat"
        },
        {
            "id": 535,
            "name": "PD-0325901"
        },
        {
            "id": 536,
            "name": "PD-0332991"
        },
        {
            "id": 554,
            "name": "PHA-665752"
        },
        {
            "id": 569,
            "name": "PLX4720"
        },
        {
            "id": 590,
            "name": "RAF265"
        },
        {
            "id": 651,
            "name": "Sorafenib"
        },
        {
            "id": 669,
            "name": "TAE684"
        },
        {
            "id": 694,
            "name": "TKI258"
        },
        {
            "id": 697,
            "name": "Topotecan"
        },
        {
            "id": 717,
            "name": "Vandetanib"
        }
    ],
    "metadata": {
        "last_page": 1,
        "page": 1,
        "per_page": 30,
        "total": 24
    }
}
```
