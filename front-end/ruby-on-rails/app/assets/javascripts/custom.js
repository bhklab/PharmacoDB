$( document ).on('turbolinks:load', function() {
	$('.slick-slider').slick({
		dots: true,
		arrows: false,
		infinite: true,
		autoplay: true,
		draggable: false,
		autoplaySpeed: 9000,
		slidesToShow: 1,
		slidesToScroll: 1
	});
});

var options = {
	useEasing : true,
	useGrouping : true,
	separator : ',',
	decimal : '.',
	prefix : '',
	suffix : ''
};

$( document ).on('turbolinks:load', function() {
	$('.hero-slider').slick({
		dots: true,
		infinite: true,
		autoplay: true,
		slidesToShow: 1,
		slidesToScroll: 1,
		autoplaySpeed: 6000,
		arrows: false
	});
});

$( document ).on('turbolinks:load', function() {
	$(function(){
		$("#search-input").typed({
		strings: ["Cell line (eg. '22rv1')", "Tissue (eg. 'endometrium')", "Drug (eg. 'paclitaxel')", "Dataset (eg. 'ccle')", "Cell line vs. Drug (eg. '22rv1 paclitaxel')", "Multiple Datasets (eg. 'ccle ctrpv2')"],
		typeSpeed: 25,
		backDelay: 1100,
		});

		$("#experiments-search-input").typed({
		strings: ["Cell line (eg. '22rv1')", "Drug (eg. 'paclitaxel')", "Cell line vs. Drug (eg. '22rv1 paclitaxel')", "Multiple Datasets (eg. 'ccle ctrpv2')"],
		typeSpeed: 25,
		backDelay: 1100,
		});
	});
});

$( document ).on('turbolinks:load', function() {
	$(function(){
		$("#experiments-search-input").typed({
		strings: ["Cell line (eg. '22rv1')", "Tissue (eg. 'endometrium')", "Drug (eg. 'paclitaxel')", "Dataset (eg. 'ccle')", "Cell line vs. Drug (eg. '22rv1 paclitaxel')", "Multiple Datasets (eg. 'ccle ctrpv2')"],
		typeSpeed: 25,
		backDelay: 1100,
		});
	});
});

/* handle smooth in-page links */

// handle links with @href started with '#' only
$(document).on('click', 'a[href^="#"]', function(e) {
  // target element id
  var id = $(this).attr('href');

  // target element
  var $id = $(id);
  if ($id.length === 0) {
    return;
  }

  // prevent standard hash navigation (avoid blinking in IE)
  e.preventDefault();

  // top position relative to the document
  var pos = $(id).offset().top;

  // animated top scrolling
  $('body, html').animate({scrollTop: pos});
});

// $( document ).on('turbolinks:load', function() {
// 	$('#experiments-search-input').typeahead({
// 		hint: false,
// 		highlight: true,
// 		minLength: 1
// 	},
// 	{
// 		limit: 5,
// 		async: true,
// 		source: function (query, processSync, processAsync) {
// 			return $.ajax({
// 				url: "/autocomplete",
// 				type: 'GET',
// 				data: {cell_query: query},
// 				dataType: 'json',
// 				success: function (json) {
// 					return processAsync(json);
// 				}
// 			});
// 		},
// 		templates: {
// 			header: '<h3 class="suggestion-header" id="cell-line-suggestion-hdr">cell line</h3>'
// 		}
// 	},
// 	{
// 		limit: 5,
// 		async: true,
// 		source: function (query, processSync, processAsync) {
// 			return $.ajax({
// 				url: "/autocomplete",
// 				type: 'GET',
// 				data: {tissue_query: query},
// 				dataType: 'json',
// 				success: function (json) {
// 					return processAsync(json);
// 				}
// 			});
// 		},
// 		templates: {
// 			header: '<h3 class="suggestion-header s-hdr-sec" id="tissue-suggestion-hdr">tissue</h3>'
// 		}
// 	},
// 	{
// 		limit: 5,
// 		async: true,
// 		source: function (query, processSync, processAsync) {
// 			return $.ajax({
// 				url: "/autocomplete",
// 				type: 'GET',
// 				data: {drug_query: query},
// 				dataType: 'json',
// 				success: function (json) {
// 					return processAsync(json);
// 				}
// 			});
// 		},
// 		templates: {
// 			header: '<h3 class="suggestion-header s-hdr-sec" id="drug-suggestion-hdr">drug</h3>'
// 		}
// 	},
// 	{
// 		limit: 5,
// 		async: true,
// 		source: function (query, processSync, processAsync) {
// 			return $.ajax({
// 				url: "/autocomplete",
// 				type: 'GET',
// 				data: {target_query: query},
// 				dataType: 'json',
// 				success: function (json) {
// 					return processAsync(json);
// 				}
// 			});
// 		},
// 		templates: {
// 			header: '<h3 class="suggestion-header s-hdr-sec" id="target-suggestion-hdr">target</h3>'
// 		}
// 	},
// 	{
// 		limit: 5,
// 		async: true,
// 		source: function (query, processSync, processAsync) {
// 			return $.ajax({
// 				url: "/autocomplete",
// 				type: 'GET',
// 				data: {dataset_query: query},
// 				dataType: 'json',
// 				success: function (json) {
// 					return processAsync(json);
// 				}
// 			});
// 		},
// 		templates: {
// 			header: '<h3 class="suggestion-header s-hdr-sec" id="dataset-suggestion-hdr">dataset</h3>'
// 		}
// 	});
// });

$( document ).on('turbolinks:load', function() {
	$('#experiments-search-input').focus();
});

$( document ).on('turbolinks:load', function() {
	$('#search-input').focus();
});

/* -------------------------------------------------------------------------- */
/* EXPLORE PAGE                                                               */
/* -------------------------------------------------------------------------- */

// Tissues
$( document ).on('click', '.annotation-tissue', function(e) {
	if ($( this ).attr('tissue-active') == 'true' ) {
		var attr = $( this ).attr('tissue-select');
		if (typeof attr !== typeof undefined && attr !== false) {
			$( this ).removeAttr('tissue-select');
			$( this ).removeAttr('style');
			$( this ).find('a').css({'color': 'black', 'font-size': '20px'});
		} else {
			$( this ).attr('tissue-select', "true");
			$( this ).css({'background-color': 'rgba(77, 77, 77, 0.5)'});
			$( this ).find('a').css({'color': '#ffffff', 'font-size': '25px'});
			$( this ).find('a i').css({'display': 'inline-block'});
		}
		var tissue_num = $("div[tissue-select='true']").length;
		var target_num = $("div[target-select='true']").length;
		$( document ).find("div[id='explore-subtitle'] h3").text("Step 1: select tissues" + (tissue_num > 0 ? " (" + tissue_num + ")" : "") + " / targets" + (target_num > 0 ? " (" + target_num + ")" : "") );
	}
});

// Cell Lines
$( document ).on('click', '.annotation-cell-line', function(e) {
	if ($( this ).attr('cell-line-active') == 'true' ) {
		var attr = $( this ).attr('cell-line-select');
		var header = $( document ).find("div[id='explore-subtitle'] h3");
		// var is_drug_selected = $("div[drug-select='true']").length > 0;
		if (typeof attr !== typeof undefined && attr !== false) {
			$( this ).removeAttr('cell-line-select');
			$( this ).removeAttr('style');
			$( this ).find('a').css({'color': 'black', 'font-size': '20px'});
			if ($( document ).find("div[drug-select='true']").length > 0) {
				header.text("Step 2: select cell line or click next to view drug");
			} else {
				header.text("Step 2: select cell line / drug");
			}
			$( '[drug-active="false"]' ).each(function( index ) {
				$( this ).attr("drug-active", "true");
				$( this ).find('a').removeAttr("style", "color: rgba(0, 0, 0, 0.5);");
			});
		} else {
			if ($("div[cell-line-select='true']").length < 1) {
				$( this ).attr('cell-line-select', "true");
				$( this ).css({'background-color': 'rgba(77, 77, 77, 0.5)'});
				$( this ).find('a').css({'color': '#ffffff', 'font-size': '25px'});
				if ($( document ).find("div[drug-select='true']").length > 0) {
					header.text("Click next to view a drug dose response curve");
				} else {
					header.text("Step 2: select drug or click next to view cell line");
				}
				// get valid drugs (for drug dose response curve) using ajax
				$.ajax({
					type:'post',
					data:{
						cid: $( this ).attr("cell-line-id")
					},
					url: 'explore/',
					dataType: 'html',
					success:function(response){
						var valid_drugs = JSON.parse(response)["drug_ids"];
						var elements = (document.getElementsByClassName('annotation-drug'));
						var drug_ids = [];
						for( var i=0; typeof(elements[i])!='undefined'; drug_ids.push(parseInt(elements[i++].getAttribute('drug-id'))));
						var diff = $(drug_ids).not(valid_drugs).get();
						for (var i = 0; i < diff.length; i++) {
							$( '[drug-id="' + diff[i] + '"]' ).each(function( index ) {
								$( this ).attr("drug-active", "false");
								$( this ).find('a').attr("style", "color: rgba(0, 0, 0, 0.5);");
							});
						}			
					}
				});
			}
		}
	}
});

// Drug Targets
$( document ).on('click', '.annotation-target', function(e) {
	if ($( this ).attr('target-active') == 'true' ) {
		var attr = $( this ).attr('target-select');
		if (typeof attr !== typeof undefined && attr !== false) {
			$( this ).removeAttr('target-select');
			$( this ).removeAttr('style');
			$( this ).find('a').css({'color': 'black', 'font-size': '20px'});
		} else {
			$( this ).attr('target-select', "true");
			$( this ).css({'background-color': 'rgba(77, 77, 77, 0.5)'});
			$( this ).find('a').css({'color': '#ffffff', 'font-size': '25px'});
		}
		var tissue_num = $("div[tissue-select='true']").length;
		var target_num = $("div[target-select='true']").length;
		$( document ).find("div[id='explore-subtitle'] h3").text("Step 1: select tissues" + (tissue_num > 0 ? " (" + tissue_num + ")" : "") + " / targets" + (target_num > 0 ? " (" + target_num + ")" : "") );	
	}
});

// Drugs
$( document ).on('click', '.annotation-drug', function(e) {
	if ($( this ).attr('drug-active') == 'true' ) {
		var attr = $( this ).attr('drug-select');
		var header = $( document ).find("div[id='explore-subtitle'] h3");
		// var is_cell_line_selected = $("div[cell-line-select='true']").length > 0;
		if (typeof attr !== typeof undefined && attr !== false) {
			$( this ).removeAttr('drug-select');
			$( this ).removeAttr('style');
			$( this ).find('a').css({'color': 'black', 'font-size': '20px'});
			if ($( document ).find("div[cell-line-select='true']").length > 0) {
				header.text("Step 2: select drug or click next to view cell line");
			} else {
				header.text("Step 2: select cell line / drug");
			}
			$( '[cell-line-active="false"]' ).each(function( index ) {
				$( this ).attr("cell-line-active", "true");
				$( this ).find('a').removeAttr("style", "color: rgba(0, 0, 0, 0.5);");
			});
		} else {
			if ($("div[drug-select='true']").length < 1) {
				$( this ).attr('drug-select', "true");
				$( this ).css({'background-color': 'rgba(77, 77, 77, 0.5)'});
				$( this ).find('a').css({'color': '#ffffff', 'font-size': '25px'});
				if ($( document ).find("div[cell-line-select='true']").length > 0) {
					header.text("Click next to view a drug dose response curve");
				} else {
					header.text("Step 2: select cell line or click next to view drug");
				}
				// get valid cell lines (for drug dose response curve) using ajax
				$.ajax({
					type:'post',
					data:{
						drug_id: $( this ).attr("drug-id")
					},
					url: 'explore/',
					dataType: 'html',
					success:function(response){
						var valid_cell_lines = JSON.parse(response)["cids"];
						var elements = (document.getElementsByClassName('annotation-cell-line'));
						var cids = [];
						for( var i=0; typeof(elements[i])!='undefined'; cids.push(parseInt(elements[i++].getAttribute('cell-line-id'))));
						var diff = $(cids).not(valid_cell_lines).get();
						for (var i = 0; i < diff.length; i++) {
							$( '[cell-line-id="' + diff[i] + '"]' ).each(function( index ) {
								$( this ).attr("cell-line-active", "false");
								$( this ).find('a').attr("style", "color: rgba(0, 0, 0, 0.5);");
							});
						}			
					}
				});
			}
		}
	}
});

$( document ).on('click', '#confirm-selection-btn', function(e) {

		tissues = "";
		cell_lines = [];
		cell_line_names = [];
		drug_targets = "";
		drugs = [];
		drug_names = [];

		$( "div[tissue-select='true']" ).each(function( index ) {
			tissues += $( this ).attr("tissue-id") + "+";
		});

		$( "div[target-select='true']" ).each(function( index ) {
			drug_targets += $( this ).attr("target-id") + "+";
		});

		$( "div[cell-line-select='true']" ).each(function( index ) {
			cell_lines.push($( this ).attr("cell-line-id"));
			cell_line_names.push($( this ).text());
		});

		$( "div[drug-select='true']" ).each(function( index ) {
			drugs.push($( this ).attr("drug-id"));
			drug_names.push($( this ).text());
		});

		var link = "http://pharmacodb.pmgenomics.ca/explore?";

		if (tissues.length > 0 || drug_targets.length > 0) {
			if (tissues.length > 0) {
					link += "tid=".concat(tissues.substring(0, tissues.length - 1));
			}
			if (drug_targets.length > 0) {
				if (tissues.length > 0) { link += "&"; }
				link += "target_id=".concat(drug_targets.substring(0, drug_targets.length - 1));
			}

		} else if (cell_lines.length > 0 || drugs.length > 0) {
			if (cell_lines.length > 0 && drugs.length > 0) {
				link = "http://pharmacodb.pmgenomics.ca/search?q=" + cell_line_names[0] + "+" + drug_names[0];
			} else {
				if (cell_lines.length > 0) {
					link = "http://pharmacodb.pmgenomics.ca/cell_lines/".concat(cell_lines[0]);
				} else {
					link = "http://pharmacodb.pmgenomics.ca/drugs/".concat(drugs[0]);
				}
			}
		} else {
			link += "select_all=true"
		}

		window.location.href = link;

});

/* -------------------------------------------------------------------------- */
/* BATCH QUERY                                                                */
/* -------------------------------------------------------------------------- */
$( document ).on('click', '#submit-batch-query-btn', function(e) {
	var batch_input = $( document ).find("div[id='batch-input']");
	var cell_lines_arr = []
	var drugs_arr = []
	// cell lines
	var cell_lines_text = $.trim(batch_input.find("textarea[id='cell-lines-ta']").val());
	if (cell_lines_text.length > 0) {
		$( document ).find(".cell-line-warning").attr("style", "display: none;");
		$.each(cell_lines_text.split(/\r?\n/), function(i, el){
			el = $.trim(el);
			if($.inArray(el, cell_lines_arr) === -1 && el.length > 0) cell_lines_arr.push(el);
		});
	} else {
		$( document ).find(".cell-line-warning").attr("style", "display: block;");
	}
	// drugs
	var drugs_text = $.trim(batch_input.find("textarea[id='drugs-ta']").val());
	if (drugs_text.length > 0) {
		$( document ).find(".drugs-warning").attr("style", "display: none;");
		$.each(drugs_text.split(/\r?\n/), function(i, el){
			el = $.trim(el);
			if($.inArray(el, drugs_arr) === -1 && el.length > 0) drugs_arr.push(el);
		});
	} else {
		$( document ).find(".drugs-warning").attr("style", "display: block;");
	}

	// create string fragments that will be sent with http post
	
	var isEmpty = cell_lines_arr.length == 0 || drugs_arr.length == 0;
	if (!isEmpty) {
		// get valid cell lines (for drug dose response curve) using ajax
		$.ajax({
			type:'post',
			data:{
				cell_names: cell_lines_arr,
				drug_names: drugs_arr
			},
			url: '/batch_query',
			dataType: 'html',
			success:function(response){
				if (response.length == "1") {
					alert('We have no results on the indicated cell-lines and compounds. Please try another query.');
				} else {
					var date = new Date();
					var dataStr = "data:text/json;charset=utf-8," + encodeURIComponent(response);
					var dlAnchorElem = document.getElementById('ajax-submit');
					dlAnchorElem.setAttribute("href", dataStr);
					dlAnchorElem.setAttribute("download", "pdb-batch-" + date.toISOString() + ".csv");
					dlAnchorElem.click();			
				}
			}
		});
	}

});

/* -------------------------------------------------------------------------- */

/* ABOUT SECTION */
/* CODE BELOW NEEDS TO BE REFACTORED */
/* -------------------------------------------------------------------------- */

function showExplore() {
	slider = document.getElementById("explore-pharmacodb");
	if (slider.value == "off" ) {
		slider.value = "on";

	} else {
		slider.value = "off";
	}
}

function show(btn, ele){
	btn_msg = document.getElementById(btn).innerText
	if (btn_msg == "see more") {
		document.getElementById(btn).innerText = "see less"
		document.getElementById(ele).style.display = 'block';
	} else {
		document.getElementById(btn).innerText = "see more"
		document.getElementById(ele).style.display = 'none';
	}
}
function on() {
    document.getElementById("overlay").style.display = "block";
    document.getElementById("home-about-container").style.display = "block"

}
function off() {
    document.getElementById("overlay").style.display = "none";
    document.getElementById("home-about-container").style.display = "none"
}

//Cookies.set('visited', Number(Cookies.get('visited')) + 1)

// window.onload = off
//if (Cookies.get('visited') == 1){on()}
/* -------------------------------------------------------------------------- */

Cookies.set('visited', Number(Cookies.get('visited')) + 1)

// window.onload = off
if (Cookies.get('visited') == 1){on()}
