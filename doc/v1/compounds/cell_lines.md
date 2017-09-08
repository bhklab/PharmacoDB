# Cell Lines Tested Against A Compound.

```
GET /compounds/{id}/cell_lines
```

## Description

This method returns a list of unique cell lines that have been treated against a compound of interest.

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

- http://pharmacodb.pmgenomics.ca/compounds

## Parameters

```
GET /compounds/{id}/cell_lines
```

| Parameter | Type | Value | Required | Default | Description |
| --- | --- | --- | --- | --- | --- |
| **page** | filter | *integer* | no | 1 | The page number for output |
| **per_page** | filter | *integer* | no | 30 | Number of items returned per page |
| **include** | input | metadata | no | - | Include meta info (eg. pagination) in body instead of headers |
| **indent** | input | *boolean* | no | false | Add indentation to response |
| **type** | input | - | no | compound_id | Define whether `id = compound_id` or `id = compound_name` |

## Output Formats

- JSON

## Examples

```
GET /compounds/{id}/cell_lines
```

- https://api.pharmacodb.com/v1/compounds/526/cell_lines
- https://api.pharmacodb.com/v1/compounds/526/cell_lines?page=2&per_page=10
- https://api.pharmacodb.com/v1/compounds/paclitaxel/cell_lines?type=name

## Output

**JSON**, using the compound **paclitaxel**, and meta info included in body.

```
{
    "data": [
        {
            "cell_line": {
                "id": 906,
                "name": "MDA-MB-453"
            },
            "datasets": [
                "CCLE",
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 9
        },
        {
            "cell_line": {
                "id": 109,
                "name": "BT-474"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 8
        },
        {
            "cell_line": {
                "id": 469,
                "name": "HCC1187"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "GRAY",
                "UHNBreast"
            ],
            "experiment_count": 8
        },
        {
            "cell_line": {
                "id": 483,
                "name": "HCC1806"
            },
            "datasets": [
                "CCLE",
                "gCSI",
                "GRAY",
                "CTRPv2",
                "UHNBreast"
            ],
            "experiment_count": 8
        },
        {
            "cell_line": {
                "id": 900,
                "name": "MDA-MB-231"
            },
            "datasets": [
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2",
                "UHNBreast"
            ],
            "experiment_count": 8
        },
        {
            "cell_line": {
                "id": 907,
                "name": "MDA-MB-468"
            },
            "datasets": [
                "CCLE",
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 8
        },
        {
            "cell_line": {
                "id": 111,
                "name": "BT-549"
            },
            "datasets": [
                "CCLE",
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2",
                "UHNBreast"
            ],
            "experiment_count": 7
        },
        {
            "cell_line": {
                "id": 895,
                "name": "MCF7"
            },
            "datasets": [
                "CCLE",
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 7
        },
        {
            "cell_line": {
                "id": 905,
                "name": "MDA-MB-436"
            },
            "datasets": [
                "CCLE",
                "gCSI",
                "GRAY",
                "FIMM"
            ],
            "experiment_count": 7
        },
        {
            "cell_line": {
                "id": 1383,
                "name": "SK-BR-3"
            },
            "datasets": [
                "CCLE",
                "gCSI",
                "GRAY",
                "FIMM"
            ],
            "experiment_count": 7
        },
        {
            "cell_line": {
                "id": 1689,
                "name": "ZR-75-1"
            },
            "datasets": [
                "CCLE",
                "GRAY",
                "FIMM",
                "CTRPv2",
                "UHNBreast"
            ],
            "experiment_count": 7
        },
        {
            "cell_line": {
                "id": 70,
                "name": "AU565"
            },
            "datasets": [
                "CCLE",
                "gCSI",
                "GRAY",
                "CTRPv2"
            ],
            "experiment_count": 6
        },
        {
            "cell_line": {
                "id": 466,
                "name": "HCC1143"
            },
            "datasets": [
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2",
                "UHNBreast"
            ],
            "experiment_count": 6
        },
        {
            "cell_line": {
                "id": 474,
                "name": "HCC1419"
            },
            "datasets": [
                "gCSI",
                "GRAY",
                "CTRPv2"
            ],
            "experiment_count": 6
        },
        {
            "cell_line": {
                "id": 489,
                "name": "HCC2185"
            },
            "datasets": [
                "GRAY"
            ],
            "experiment_count": 6
        },
        {
            "cell_line": {
                "id": 614,
                "name": "Hs-578-T"
            },
            "datasets": [
                "CCLE",
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 6
        },
        {
            "cell_line": {
                "id": 898,
                "name": "MDA-MB-157"
            },
            "datasets": [
                "CCLE",
                "gCSI",
                "GRAY",
                "CTRPv2",
                "UHNBreast"
            ],
            "experiment_count": 6
        },
        {
            "cell_line": {
                "id": 1559,
                "name": "T47D"
            },
            "datasets": [
                "CCLE",
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 6
        },
        {
            "cell_line": {
                "id": 1635,
                "name": "UACC-812"
            },
            "datasets": [
                "CCLE",
                "GRAY"
            ],
            "experiment_count": 6
        },
        {
            "cell_line": {
                "id": 1636,
                "name": "UACC-893"
            },
            "datasets": [
                "GRAY"
            ],
            "experiment_count": 6
        },
        {
            "cell_line": {
                "id": 48,
                "name": "A549"
            },
            "datasets": [
                "CCLE",
                "gCSI",
                "CTRPv2"
            ],
            "experiment_count": 5
        },
        {
            "cell_line": {
                "id": 337,
                "name": "ES-2"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 5
        },
        {
            "cell_line": {
                "id": 473,
                "name": "HCC1395"
            },
            "datasets": [
                "CCLE",
                "GRAY",
                "CTRPv2"
            ],
            "experiment_count": 5
        },
        {
            "cell_line": {
                "id": 485,
                "name": "HCC1937"
            },
            "datasets": [
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2",
                "UHNBreast"
            ],
            "experiment_count": 5
        },
        {
            "cell_line": {
                "id": 486,
                "name": "HCC1954"
            },
            "datasets": [
                "CCLE",
                "gCSI",
                "GRAY"
            ],
            "experiment_count": 5
        },
        {
            "cell_line": {
                "id": 654,
                "name": "IGROV-1"
            },
            "datasets": [
                "CCLE",
                "gCSI",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 5
        },
        {
            "cell_line": {
                "id": 857,
                "name": "LOXIMVI"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "CTRPv2"
            ],
            "experiment_count": 5
        },
        {
            "cell_line": {
                "id": 897,
                "name": "MDA-MB-134-VI"
            },
            "datasets": [
                "GDSC1000",
                "GRAY",
                "UHNBreast"
            ],
            "experiment_count": 5
        },
        {
            "cell_line": {
                "id": 902,
                "name": "MDA-MB-361"
            },
            "datasets": [
                "GRAY",
                "CTRPv2"
            ],
            "experiment_count": 5
        },
        {
            "cell_line": {
                "id": 1043,
                "name": "NCI-H1869"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "CTRPv2"
            ],
            "experiment_count": 5
        }
    ],
    "metadata": {
        "last_page": 40,
        "page": 1,
        "per_page": 30,
        "total": 1174
    }
}
```
