# Tissues

```
GET /compounds/{id}/tissues
```

## Description

This method returns a list of unique tissues that have been treated against a compound of interest.

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
GET /compounds/{id}/tissues
```

| Parameter | Type | Value | Required | Default | Description |
| --- | --- | --- | --- | --- | --- |
| **page** | filter | *integer* | no | 1 | The page number for output |
| **per_page** | filter | *integer* | no | 30 | Number of items returned per page |
| **include** | input | metadata | no | - | Include meta info (eg. pagination) in body instead of headers |
| **indent** | input | *boolean* | no | false | Add indentation to response |
| **type** | input | - | no | tissue_id | Define whether `id = compound_id` or `id = compound_name` |

## Output Formats

- JSON

## Examples

```
GET /compounds/{id}/tissues
```

- https://api.pharmacodb.com/v1/compounds/526/tissues
- https://api.pharmacodb.com/v1/compounds/526/tissues?page=2&per_page=10
- https://api.pharmacodb.com/v1/compounds/paclitaxel/tissues?type=name

## Output

**JSON**, using the compound **paclitaxel**, and meta info included in body.

```
{
    "data": [
        {
            "tissue": {
                "id": 13,
                "name": "haematopoietic_and_lymphoid_tissue"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 458
        },
        {
            "tissue": {
                "id": 18,
                "name": "lung"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "CTRPv2"
            ],
            "experiment_count": 382
        },
        {
            "tissue": {
                "id": 7,
                "name": "breast"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "GRAY",
                "FIMM",
                "CTRPv2",
                "UHNBreast"
            ],
            "experiment_count": 297
        },
        {
            "tissue": {
                "id": 31,
                "name": "skin"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "CTRPv2"
            ],
            "experiment_count": 130
        },
        {
            "tissue": {
                "id": 9,
                "name": "central_nervous_system"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "CTRPv2"
            ],
            "experiment_count": 118
        },
        {
            "tissue": {
                "id": 16,
                "name": "large_intestine"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "CTRPv2"
            ],
            "experiment_count": 113
        },
        {
            "tissue": {
                "id": 24,
                "name": "ovary"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 109
        },
        {
            "tissue": {
                "id": 25,
                "name": "pancreas"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "CTRPv2"
            ],
            "experiment_count": 99
        },
        {
            "tissue": {
                "id": 34,
                "name": "stomach"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "CTRPv2"
            ],
            "experiment_count": 75
        },
        {
            "tissue": {
                "id": 2,
                "name": "autonomic_ganglia"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "CTRPv2"
            ],
            "experiment_count": 65
        },
        {
            "tissue": {
                "id": 11,
                "name": "endometrium"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "CTRPv2"
            ],
            "experiment_count": 64
        },
        {
            "tissue": {
                "id": 21,
                "name": "oesophagus"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "CTRPv2"
            ],
            "experiment_count": 60
        },
        {
            "tissue": {
                "id": 17,
                "name": "liver"
            },
            "datasets": [
                "CCLE",
                "gCSI",
                "CTRPv2"
            ],
            "experiment_count": 58
        },
        {
            "tissue": {
                "id": 5,
                "name": "bone"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "CTRPv2"
            ],
            "experiment_count": 52
        },
        {
            "tissue": {
                "id": 37,
                "name": "upper_aerodigestive_tract"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "CTRPv2"
            ],
            "experiment_count": 49
        },
        {
            "tissue": {
                "id": 15,
                "name": "kidney"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "CTRPv2"
            ],
            "experiment_count": 48
        },
        {
            "tissue": {
                "id": 38,
                "name": "urinary_tract"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "CTRPv2"
            ],
            "experiment_count": 47
        },
        {
            "tissue": {
                "id": 41,
                "name": "N/A"
            },
            "datasets": [
                "CTRPv2"
            ],
            "experiment_count": 45
        },
        {
            "tissue": {
                "id": 33,
                "name": "soft_tissue"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "CTRPv2"
            ],
            "experiment_count": 44
        },
        {
            "tissue": {
                "id": 36,
                "name": "thyroid"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "CTRPv2"
            ],
            "experiment_count": 21
        },
        {
            "tissue": {
                "id": 27,
                "name": "pleura"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "CTRPv2"
            ],
            "experiment_count": 20
        },
        {
            "tissue": {
                "id": 28,
                "name": "prostate"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "FIMM",
                "CTRPv2"
            ],
            "experiment_count": 15
        },
        {
            "tissue": {
                "id": 3,
                "name": "biliary_tract"
            },
            "datasets": [
                "CCLE",
                "GDSC1000",
                "gCSI",
                "CTRPv2"
            ],
            "experiment_count": 11
        },
        {
            "tissue": {
                "id": 10,
                "name": "cervix"
            },
            "datasets": [
                "GDSC1000",
                "gCSI",
                "CTRPv2"
            ],
            "experiment_count": 11
        },
        {
            "tissue": {
                "id": 35,
                "name": "testis"
            },
            "datasets": [
                "GDSC1000"
            ],
            "experiment_count": 4
        },
        {
            "tissue": {
                "id": 40,
                "name": "vulva"
            },
            "datasets": [
                "GDSC1000"
            ],
            "experiment_count": 2
        },
        {
            "tissue": {
                "id": 12,
                "name": "gastrointestinal_tract_(site_indeterminate)"
            },
            "datasets": [
                "GDSC1000"
            ],
            "experiment_count": 1
        },
        {
            "tissue": {
                "id": 20,
                "name": "neuroblastoma"
            },
            "datasets": [
                "GDSC1000"
            ],
            "experiment_count": 1
        },
        {
            "tissue": {
                "id": 29,
                "name": "salivary_gland"
            },
            "datasets": [
                "CTRPv2"
            ],
            "experiment_count": 1
        },
        {
            "tissue": {
                "id": 32,
                "name": "small_intestine"
            },
            "datasets": [
                "GDSC1000"
            ],
            "experiment_count": 1
        }
    ],
    "metadata": {
        "last_page": 1,
        "page": 1,
        "per_page": 30,
        "total": 30
    }
}
```
