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

$( document ).on('turbolinks:load', function() {
	$('#search-input').typeahead({
		hint: false,
		highlight: true,
		minLength: 1
	},
	{
		limit: 5,
		async: true,
		source: function (query, processSync, processAsync) {
			return $.ajax({
				url: "/autocomplete",
				type: 'GET',
				data: {cell_query: query},
				dataType: 'json',
				success: function (json) {
					return processAsync(json);
				}
			});
		},
		templates: {
			header: '<h3 class="suggestion-header" id="cell-line-suggestion-hdr">cell line</h3>'
		}
	},
	{
		limit: 5,
		async: true,
		source: function (query, processSync, processAsync) {
			return $.ajax({
				url: "/autocomplete",
				type: 'GET',
				data: {tissue_query: query},
				dataType: 'json',
				success: function (json) {
					return processAsync(json);
				}
			});
		},
		templates: {
			header: '<h3 class="suggestion-header s-hdr-sec" id="tissue-suggestion-hdr">tissue</h3>'
		}
	},
	{
		limit: 5,
		async: true,
		source: function (query, processSync, processAsync) {
			return $.ajax({
				url: "/autocomplete",
				type: 'GET',
				data: {drug_query: query},
				dataType: 'json',
				success: function (json) {
					return processAsync(json);
				}
			});
		},
		templates: {
			header: '<h3 class="suggestion-header s-hdr-sec" id="drug-suggestion-hdr">drug</h3>'
		}
	},
	{
		limit: 5,
		async: true,
		source: function (query, processSync, processAsync) {
			return $.ajax({
				url: "/autocomplete",
				type: 'GET',
				data: {target_query: query},
				dataType: 'json',
				success: function (json) {
					return processAsync(json);
				}
			});
		},
		templates: {
			header: '<h3 class="suggestion-header s-hdr-sec" id="target-suggestion-hdr">target</h3>'
		}
	},
	{
		limit: 5,
		async: true,
		source: function (query, processSync, processAsync) {
			return $.ajax({
				url: "/autocomplete",
				type: 'GET',
				data: {dataset_query: query},
				dataType: 'json',
				success: function (json) {
					return processAsync(json);
				}
			});
		},
		templates: {
			header: '<h3 class="suggestion-header s-hdr-sec" id="dataset-suggestion-hdr">dataset</h3>'
		}
	},
	{
			limit: 5,
			async: true,
			source: function (query, processSync, processAsync) {
				return $.ajax({
					url: "/autocomplete",
					type: 'GET',
					data: {other_query: query},
					dataType: 'json',
					success: function (json) {
						return processAsync(json);
					}
				});
			},
			templates: {
				header: '<h3 class="suggestion-header s-hdr-sec" id="other-suggestion-hdr">No Results</h3>'
			}
		});
});

$( document ).on('turbolinks:load', function() {
	$('#experiments-search-input').typeahead({
		hint: false,
		highlight: true,
		minLength: 1
	},
	{
		limit: 5,
		async: true,
		source: function (query, processSync, processAsync) {
			return $.ajax({
				url: "/autocomplete",
				type: 'GET',
				data: {cell_query: query},
				dataType: 'json',
				success: function (json) {
					return processAsync(json);
				}
			});
		},
		templates: {
			header: '<h3 class="suggestion-header" id="cell-line-suggestion-hdr">cell line</h3>'
		}
	},
	{
		limit: 5,
		async: true,
		source: function (query, processSync, processAsync) {
			return $.ajax({
				url: "/autocomplete",
				type: 'GET',
				data: {tissue_query: query},
				dataType: 'json',
				success: function (json) {
					return processAsync(json);
				}
			});
		},
		templates: {
			header: '<h3 class="suggestion-header s-hdr-sec" id="tissue-suggestion-hdr">tissue</h3>'
		}
	},
	{
		limit: 5,
		async: true,
		source: function (query, processSync, processAsync) {
			return $.ajax({
				url: "/autocomplete",
				type: 'GET',
				data: {drug_query: query},
				dataType: 'json',
				success: function (json) {
					return processAsync(json);
				}
			});
		},
		templates: {
			header: '<h3 class="suggestion-header s-hdr-sec" id="drug-suggestion-hdr">drug</h3>'
		}
	},
	{
		limit: 5,
		async: true,
		source: function (query, processSync, processAsync) {
			return $.ajax({
				url: "/autocomplete",
				type: 'GET',
				data: {target_query: query},
				dataType: 'json',
				success: function (json) {
					return processAsync(json);
				}
			});
		},
		templates: {
			header: '<h3 class="suggestion-header s-hdr-sec" id="target-suggestion-hdr">target</h3>'
		}
	},
	{
		limit: 5,
		async: true,
		source: function (query, processSync, processAsync) {
			return $.ajax({
				url: "/autocomplete",
				type: 'GET',
				data: {dataset_query: query},
				dataType: 'json',
				success: function (json) {
					return processAsync(json);
				}
			});
		},
		templates: {
			header: '<h3 class="suggestion-header s-hdr-sec" id="dataset-suggestion-hdr">dataset</h3>'
		}
	});
});

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
$( document ).on('click', '#tissue-row', function(e) {
	if ($( this ).attr('tissue-active') == 'true' ) {
		$( this ).attr('tissue-select', "true");
		$( this ).css({'background-color': 'rgba(250, 255, 0, 1.0)'});
		$( this ).find('td a').css({'color': '#ff0084', 'font-size': '25px'});
	}
});

// Cell Lines
$( document ).on('click', '#cell-line-row', function(e) {
	if ($( this ).attr('cell-line-active') == 'true' ) {
		$( this ).attr('cell-line-select', "true");
		$( this ).css({'background-color': 'rgba(250, 255, 0, 1.0)'});
		$( this ).find('td a').css({'color': '#ff0084', 'font-size': '25px'});
	}
});

// Drug Targets
$( document ).on('click', '#target-row', function(e) {
	if ($( this ).attr('target-active') == 'true' ) {
		$( this ).attr('target-select', "true");
		$( this ).css({'background-color': 'rgba(250, 255, 0, 1.0)'});
		$( this ).find('td a').css({'color': '#ff0084', 'font-size': '25px'});
	}
});

// Drugs
// $( document ).on('click', '#drug-row', function(e) {
// 	if ($( this ).attr('drug-active') == 'true' ) {
// 		$( this ).attr('drug-select', "true");
// 		$( this ).css({'background-color': 'rgba(250, 255, 0, 1.0)'});
// 		$( this ).find('td a').css({'color': '#ff0084', 'font-size': '25px'});
// 	}
// });

$( document ).on('click', '.confirm-selection-btn', function(e) {

		tissues = "";
		cell_lines= "";
		drug_targets = "";
		drugs = "";

		$( "tr[tissue-select='true']" ).each(function( index ) {
			tissues += $( this ).attr("tissue-id") + "+";
		});

		$( "tr[cell-line-select='true']" ).each(function( index ) {
			cell_lines += $( this ).attr("cell-line-id") + "+";
		});

		if ( cell_lines.length > 0 ) {
			tissues = "";
			$( "tr[id='tissue-row']" ).each(function( index ) {
				tissues += $( this ).attr("tissue-id") + "+";
			});
		}

		$( "tr[target-select='true']" ).each(function( index ) {
			drug_targets += $( this ).attr("target-id") + "+";
		});

		// $( "tr[drug-select='true']" ).each(function( index ) {
		// 	drugs += $( this ).attr("drug-id") + "+";
		// });

		$.ajax({
		    type:'post',
			data:{
			    tid: tissues.substring(0, tissues.length - 1),
			    // cid: cell_lines.substring(0, cell_lines.length - 1),
			    // drug_id: drugs.substring(0, drugs.length - 1),
			    target_id: drug_targets.substring(0, drug_targets.length - 1),
					"authenticity_token": "<%= form_authenticity_token %>"
			},
		    url: 'explore/',
		    dataType: 'html',
		    success:function(response){
		        $('body').html(response);
		    }
		});

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
