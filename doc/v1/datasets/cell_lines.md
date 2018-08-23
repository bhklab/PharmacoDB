# Cell Lines Tested In A Dataset

```
GET /datasets/{id}/cell_lines
```

## Description

This method returns a list of unique cell lines that have been tested in a dataset of interest.

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
GET /datasets/{id}/cell_lines
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
GET /datasets/{id}/cell_lines
```

- http://api.pharmacodb.ca/v1/datasets/1/cell_lines
- http://api.pharmacodb.ca/v1/datasets/1/cell_lines?page=2&per_page=10
- http://api.pharmacodb.ca/v1/datasets/ccle/cell_lines?type=name

## Output

**JSON**, using the dataset **CCLE**, and meta info included in body.

```
{
    "data": [
        {
            "id": 2,
            "name": "1321N1"
        },
        {
            "id": 12,
            "name": "22RV1"
        },
        {
            "id": 16,
            "name": "42-MG-BA"
        },
        {
            "id": 20,
            "name": "5637"
        },
        {
            "id": 24,
            "name": "639-V"
        },
        {
            "id": 26,
            "name": "697"
        },
        {
            "id": 27,
            "name": "769-P"
        },
        {
            "id": 28,
            "name": "786-0"
        },
        {
            "id": 29,
            "name": "8-MG-BA"
        },
        {
            "id": 30,
            "name": "8305C"
        },
        {
            "id": 31,
            "name": "8505C"
        },
        {
            "id": 36,
            "name": "A172"
        },
        {
            "id": 37,
            "name": "A204"
        },
        {
            "id": 38,
            "name": "A2058"
        },
        {
            "id": 39,
            "name": "A253"
        },
        {
            "id": 40,
            "name": "A2780"
        },
        {
            "id": 42,
            "name": "A375"
        },
        {
            "id": 48,
            "name": "A549"
        },
        {
            "id": 49,
            "name": "A673"
        },
        {
            "id": 54,
            "name": "ACHN"
        },
        {
            "id": 59,
            "name": "ALL-SIL"
        },
        {
            "id": 62,
            "name": "AMO-1"
        },
        {
            "id": 63,
            "name": "AN3-CA"
        },
        {
            "id": 67,
            "name": "AsPC-1"
        },
        {
            "id": 70,
            "name": "AU565"
        },
        {
            "id": 72,
            "name": "AZ-521"
        },
        {
            "id": 82,
            "name": "BCPAP"
        },
        {
            "id": 83,
            "name": "BDCM"
        },
        {
            "id": 89,
            "name": "BFTC-909"
        },
        {
            "id": 90,
            "name": "BGC-823"
        }
    ],
    "metadata": {
        "last_page": 17,
        "page": 1,
        "per_page": 30,
        "total": 504
    }
}
```
