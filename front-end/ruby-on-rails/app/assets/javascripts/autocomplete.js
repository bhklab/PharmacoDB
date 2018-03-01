function mod(n, m) {
        return ((n % m) + m) % m;
}


// // https://stackoverflow.com/questions/3730510/javascript-sort-array-and-return-an-array-of-indicies-that-indicates-the-positi
// function sortWithIndeces(toSort) {
//   for (var i = 0; i < toSort.length; i++) {
//     toSort[i] = [toSort[i], i];
//   }
//   toSort.sort(function(left, right) {
//     return left[0] < right[0] ? -1 : 1;
//   });
//   toSort.sortIndices = [];
//   for (var j = 0; j < toSort.length; j++) {
//     toSort.sortIndices.push(toSort[j][1]);
//     toSort[j] = toSort[j][0];
//   }
//   return toSort;
// }




$( document ).on('turbolinks:load', function() {
    SC = $('#search-form').children("")
    $('#search-form').prepend('<span id="autospan" style="position: relative; display: inline-block; width: calc(100% - 5px)"> </span>')
    $('#search-form').children().first().append(SC)
    $("#autospan").append('<div id="autores" class="tt-menu" style="position: absolute; top: 100%; left: 0px; z-index: 100; display: none;"> \
     </div>')
    $("#search-form").on('focusin', function(){
        if ($('#search-input').val().length > 1){
            $("#autores").show()
        }
    })
    $("#search-form").on("focusout", function(){
        // XXX::HACK!
        window.setTimeout(function() { $("#autores").hide() }, 100);
    })


    $('#search-input').attr('autocomplete', 'off')

    var oldVal1 = $("#search-input").val()

    $("#search-input").on("keyup", function(e){
        // $("#autores").html("<h1> HI</h1>")
        if ($("#search-input").val() == oldVal1) {
            return
        }
        if ((e.keyCode >= 37 && e.keyCode <= 40) || e.keyCode==13){ 
            return 
        }
        var query = $(this).val()
        oldVal1 = query

        if (query.length < 2){
            $("#autores").hide()
            return
        } else {
            $("#autores").show()
        }
        // console.log(query)
        $.ajax({url: "/autocomplete",
           type: 'GET',
           data: {query: query},
           dataType: 'json',
           success: function (json) {
               updateSelection = updateWithRes(json);

               $("#search-input").off("keydown").on("keydown", function(e){
                if (e.keyCode == 40 || e.keyCode == 38){
                    e.preventDefault()
                    updateSelection(e)
                    // console.log(e)
                }
            })
           }
       })
    })

    

    

    // TODO: make only current input STRONG
    // TODO: implement sorting by creating html substrings and then appending things in right order before setting things
    // also do this for the optionArray
    // fix length of optionArray 
    function updateWithRes(json, si) {

        var html = "";
        var optionArray = [[$("#search-input").val()]];
        var indexArray = [];
        var j = 0;

        var htmlCell = "";
        var i ;
        var optionCellArray = [];
        var indexCellArray = [];

        if(json.cell.length){
            htmlCell += '<div id="cell-suggestion"> <h3 class="suggestion-header" id="cell-line-suggestion-hdr">Cell Line</h3>' 
        for(i = 0; i < Math.min(json.cell.length,5); i++){
            htmlCell += '<div class="tt-suggestion tt-selectable" ' + 'id="'+ j.toString() +'-auto"><strong class="tt-highlight">' + json.cell[i][0] + '</strong></div>'
            indexCellArray.push(j)
            j += 1;
        }
        htmlCell += '</div>'
        if (i > 0){
            optionCellArray = json.cell.slice(0,i)
        }
        }


        var htmlDrug = "";
        var i ;
        var optionDrugArray = [];
        var indexDrugArray = [];
        if(json.drug.length){
            htmlDrug += '<div id="drug-suggestion"> <h3 class="suggestion-header s-hdr-sec" id="drug-suggestion-hdr">Compound</h3>' 
            for(i = 0; i < Math.min(json.drug.length,5); i++){
                htmlDrug += '<div class="tt-suggestion tt-selectable" ' + 'id="'+ j.toString() +'-auto"><strong class="tt-highlight">' + json.drug[i][0] + '</strong></div>'
                indexDrugArray.push(j)
                j += 1;
            }
            htmlDrug += "</div>"
            if (i > 0) {
                optionDrugArray = json.drug.slice(0,i)
            }
        }

        var htmlTissue = "";
        var i ;
        var optionTissueArray = [];
        var indexTissueArray = [];

        if(json.tissue.length){
            htmlTissue += '<div id="tissue-suggestion"> <h3 class="suggestion-header s-hdr-sec" id="tissue-suggestion-hdr">Tissue</h3>' 
        for(i = 0; i < Math.min(json.tissue.length,5); i++){
            htmlTissue += '<div class="tt-suggestion tt-selectable" ' + 'id="'+ j.toString() +'-auto"><strong class="tt-highlight">' + json.tissue[i][0] + '</strong></div>'
            indexTissueArray.push(j)
            j += 1;
        }
        htmlTissue += '</div>'
        if (i > 0){
            optionTissueArray = json.tissue.slice(0,i)
        }

        }

        var htmlGene = "";
        var i ;
        var optionGeneArray = [];
        var indexGeneArray = [];

        if(json.gene.length){
            htmlGene += '<div id="gene-suggestion"> <h3 class="suggestion-header s-hdr-sec" id="gene-suggestion-hdr">Gene</h3>' 

            for(i = 0; i < Math.min(json.gene.length,5); i++){
                htmlGene += '<div class="tt-suggestion tt-selectable" ' + 'id="'+ j.toString() +'-auto"><strong class="tt-highlight">' + json.gene[i][0] + '</strong></div>'
                indexGeneArray.push(j)                
                j += 1;
            }
            htmlGene += '</div>'
            if (i > 0){
                optionGeneArray= json.gene.slice(0,i)
            }

        }


        var htmlDataset = "";
        var i ;
        var optionDatasetArray = [];
        var indexDatasetArray = [];

        if(json.dataset.length){
            htmlDataset += '<div id="dataset-suggestion" > <h3 class="suggestion-header s-hdr-sec" id="dataset-suggestion-hdr">Dataset</h3>' 
            for(i = 0; i < Math.min(json.dataset.length,5); i++){
                htmlDataset += '<div class="tt-suggestion tt-selectable" ' + 'id="'+ j.toString() +'-auto"><strong class="tt-highlight">' + json.dataset[i][0] + '</strong></div>'
                indexDatasetArray.push(j)
                j += 1;
            }
            htmlDataset += '</div>'
            if (i > 0){
                optionDatasetArray = json.dataset.slice(0,i)
            }

        }



        var htmlCell_Drug = "";
        var i;
        var optionCell_DrugArray = [];
        var indexCell_DrugArray = [];

        for (var i = 0; i < json.cell_drug.length; i++){
            json.cell_drug[i]
        }


        if(json.cell_drug.length){
            htmlCell_Drug += '<div id="DDRC-suggestion"> <h3 class="suggestion-header s-hdr-sec" id="cell_drug-suggestion-hdr">Drug Dose Response Curve</h3>' 

            for(i = 0; i < Math.min(json.cell_drug.length,5); i++){
                htmlCell_Drug += '<div class="tt-suggestion tt-selectable" ' + 'id="'+ j.toString() +'-auto"><strong class="tt-highlight">' + json.cell_drug[i][0] + '</strong></div>'
                indexCell_DrugArray.push(j)
                j += 1;
            }
            htmlCell_Drug += '</div>'
            if( i > 0){
                optionCell_DrugArray = json.cell_drug.slice(0,i)
            }

        }


        var htmlDataset_Int = "";
        var i;
        var optionDataset_IntArray = [];
        var indexDataset_IntArray = [];

        if(json.dataset_int.length){
            htmlDataset_Int += '<div id="Data-int-suggestion"> <h3 class="suggestion-header s-hdr-sec" id="dataset_int-suggestion-hdr">Dataset Intersection</h3>' 

            for(i = 0; i < Math.min(json.dataset_int.length,5); i++){
                htmlDataset_Int += '<div class="tt-suggestion tt-selectable" ' + 'id="'+ j.toString() +'-auto"><strong class="tt-highlight">' + json.dataset_int[i][0] + '</strong></div>'
                indexDataset_IntArray.push(j)
                j += 1;
            }
            htmlDataset_Int += '</div>'
            if( i > 0){
                optionDataset_IntArray = json.dataset_int.slice(0,i)
            }

        }


        var optionArrays = [optionCellArray, optionDrugArray, optionTissueArray, optionGeneArray, optionDatasetArray, optionCell_DrugArray, optionDataset_IntArray]
        var htmls = [htmlCell, htmlDrug, htmlTissue, htmlGene, htmlDataset, htmlCell_Drug, htmlDataset_Int]
        var indexArrays = [indexCellArray, indexDrugArray, indexTissueArray, indexGeneArray, indexDatasetArray, indexCell_DrugArray, indexDataset_IntArray]

        var toSort = [json.cell.length ? json.cell[0][1] : 0,
                      json.drug.length ? json.drug[0][1] : 0,
                      json.tissue.length ? json.tissue[0][1] : 0,
                      json.gene.length ? json.gene[0][1] : 0,
                      json.dataset.length ? json.dataset[0][1] : 0,
                      json.cell_drug.length ? json.cell_drug[0][1] : 0,
                      json.dataset_int.length ? json.dataset_int[0][1] : 0]


        var indexedTest = toSort.map(function(e,i){return {ind: i, val: e}});
        // sort index/value couples, based on values
        indexedTest.sort(function(x, y){return x.val > y.val ? -1 : x.val == y.val ? 0 : 1});
        // make list keeping only indices
        var sortOrder = indexedTest.map(function(e){return e.ind});

        for (var iii = 0; iii < sortOrder.length; iii++){
            optionArray = optionArray.concat(optionArrays[sortOrder[iii]])
            html += htmls[sortOrder[iii]]
            indexArray = indexArray.concat(indexArrays[sortOrder[iii]])
        }
        var ii = 0;
        function updateSelection(e){
            if (ii != 0){
                $("#"+ (indexArray[ii - 1]).toString() + "-auto").toggleClass("tt-selected")
            }

            // console.log(optionArray)
            if(e.keyCode == 40){
                ii = mod((ii + 1), (optionArray.length))
            }
            if(e.keyCode == 38){
                ii = mod((ii - 1), (optionArray.length))
            }
            $("#search-input").val(optionArray[ii][0])
            if (ii != 0){
                $("#"+ (indexArray[ii - 1]).toString() + "-auto").scrollintoview().toggleClass("tt-selected")
            }
            
        }
        $("#autores").html(html)

        addClickHandler = function(i){
            $("#"+ (indexArray[i-1]).toString() + "-auto").off("click").on("click", function(){
                $("#search-input").val(optionArray[i][0])
            })
        }


        for (var i = 1; i <= j; i++){
            var value = optionArray[i][0]
            addClickHandler(i)
        }
        return updateSelection
    }

});