# Cell Lines

```
GET /tissues/{id}/cell_lines
```

## Description

This method returns a list of unique cell lines that belong to a tissue type of interest.

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

- http://pharmacodb.pmgenomics.ca/tissues

## Parameters

```
GET /tissues/{id}/cell_lines
```

| Parameter | Type | Value | Required | Default | Description |
| --- | --- | --- | --- | --- | --- |
| **page** | filter | *integer* | no | 1 | The page number for output |
| **per_page** | filter | *integer* | no | 30 | Number of items returned per page |
| **include** | input | metadata | no | - | Include meta info (eg. pagination) in body instead of headers |
| **indent** | input | *boolean* | no | false | Add indentation to response |
| **type** | input | - | no | tissue_id | Define whether `id = tissue_id` or `id = tissue_name` |

## Output Formats

- JSON

## Examples

```
GET /tissues/{id}/cell_lines
```

- https://api.pharmacodb.com/v1/cell_lines/7/cell_lines
- https://api.pharmacodb.com/v1/cell_lines/7/cell_lines?page=2&per_page=10
- https://api.pharmacodb.com/v1/cell_lines/breast/cell_lines?type=name

## Output

**JSON**, using the tissue **breast**, and meta info included in body.

```
{
    "data": [
        {
            "id": 4,
            "name": "184A1"
        },
        {
            "id": 5,
            "name": "184B5"
        },
        {
            "id": 8,
            "name": "21MT1"
        },
        {
            "id": 9,
            "name": "21MT2"
        },
        {
            "id": 10,
            "name": "21NT"
        },
        {
            "id": 11,
            "name": "21PT"
        },
        {
            "id": 22,
            "name": "600MPE"
        },
        {
            "id": 70,
            "name": "AU565"
        },
        {
            "id": 71,
            "name": "AU655"
        },
        {
            "id": 108,
            "name": "BT-20"
        },
        {
            "id": 109,
            "name": "BT-474"
        },
        {
            "id": 110,
            "name": "BT-483"
        },
        {
            "id": 111,
            "name": "BT-549"
        },
        {
            "id": 165,
            "name": "CAL-120"
        },
        {
            "id": 167,
            "name": "CAL-148"
        },
        {
            "id": 172,
            "name": "CAL-51"
        },
        {
            "id": 177,
            "name": "CAL-85-1"
        },
        {
            "id": 181,
            "name": "CAMA-1"
        },
        {
            "id": 237,
            "name": "COLO-824"
        },
        {
            "id": 310,
            "name": "DU-4475"
        },
        {
            "id": 321,
            "name": "EFM-19"
        },
        {
            "id": 322,
            "name": "EFM-192A"
        },
        {
            "id": 323,
            "name": "EFM-192B"
        },
        {
            "id": 324,
            "name": "EFM-192C"
        },
        {
            "id": 349,
            "name": "EVSA-T"
        },
        {
            "id": 442,
            "name": "HARA"
        },
        {
            "id": 443,
            "name": "HBL-100"
        },
        {
            "id": 464,
            "name": "HCC1008"
        },
        {
            "id": 466,
            "name": "HCC1143"
        },
        {
            "id": 467,
            "name": "HCC1143BL"
        }
    ],
    "metadata": {
        "last_page": 4,
        "page": 1,
        "per_page": 30,
        "total": 116
    }
}
```
