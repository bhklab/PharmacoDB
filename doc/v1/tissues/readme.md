# Tissues

```
GET /tissues
```

## Description

This method returns a list of tissues.

## Summary

| Name | Value |
| --- | --- |
| **Request Protocol** | GET |
| **Requires API Key** | No |
| **Cache Time** | 0 seconds |

## Notes

- Meta information is included in response headers by default. Use `include` parameter to add info to response body.

## Sources

- http://pharmacodb.ca/tissues

## Parameters

```
GET /tissues
```

| Parameter | Type | Value | Required | Default | Description |
| --- | --- | --- | --- | --- | --- |
| **page** | filter | *integer* | no | 1 | The page number for output |
| **per_page** | filter | *integer* | no | 30 | Number of items returned per page |
| **all** | filter | *boolean* | no | false | Return all items in resource |
| **include** | input | metadata | no | - | Include meta info (eg. pagination) in body instead of headers |
| **indent** | input | *boolean* | no | false | Add indentation to response |

## Output Formats

- JSON

## Examples

```
GET /tissues
```

- http://api.pharmacodb.ca/v1/tissues

## Output

**JSON**, with metadata included in body.

```
{
    "data": [
        {
            "id": 1,
            "name": "adrenal_gland"
        },
        {
            "id": 2,
            "name": "autonomic_ganglia"
        },
        {
            "id": 3,
            "name": "biliary_tract"
        },
        {
            "id": 4,
            "name": "blood"
        },
        {
            "id": 5,
            "name": "bone"
        },
        {
            "id": 6,
            "name": "brain"
        },
        {
            "id": 7,
            "name": "breast"
        },
        {
            "id": 8,
            "name": "cartilage"
        },
        {
            "id": 9,
            "name": "central_nervous_system"
        },
        {
            "id": 10,
            "name": "cervix"
        },
        {
            "id": 11,
            "name": "endometrium"
        },
        {
            "id": 12,
            "name": "gastrointestinal_tract_(site_indeterminate)"
        },
        {
            "id": 13,
            "name": "haematopoietic_and_lymphoid_tissue"
        },
        {
            "id": 14,
            "name": "head_and_neck"
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
            "id": 19,
            "name": "misc"
        },
        {
            "id": 20,
            "name": "neuroblastoma"
        },
        {
            "id": 21,
            "name": "oesophagus"
        },
        {
            "id": 22,
            "name": "oral_cavity"
        },
        {
            "id": 23,
            "name": "other"
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
            "id": 26,
            "name": "placenta"
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
            "id": 29,
            "name": "salivary_gland"
        },
        {
            "id": 30,
            "name": "skeletal_muscle"
        }
    ],
    "metadata": {
        "last_page": 2,
        "page": 1,
        "per_page": 30,
        "total": 41
    }
}
```
