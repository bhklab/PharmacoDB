function updateDimensions(winWidth, margin) {
  width = winWidth - margin.left - margin.right;
  height = 500 - margin.top - margin.bottom;
}

function wrap(text, width) {
  text.each(function() {
    var text = d3.select(this),
        words = text.text().split(/\s+/).reverse(),
        word,
        line = [],
        lineNumber = 0,
        lineHeight = 1.1, // ems
        y = text.attr("y"),
        dy = parseFloat(text.attr("dy")),
        tspan = text.text(null).append("tspan").attr("x", 0).attr("y", y).attr("dy", dy + "em");
    while (word = words.pop()) {
      line.push(word);
      tspan.text(line.join(" "));
      if (tspan.node().getComputedTextLength() > width) {
        line.pop();
        tspan.text(line.join(" "));
        line = [word];
        tspan = text.append("tspan").attr("x", 0).attr("y", y).attr("dy", ++lineNumber * lineHeight + dy + "em").text(word);
      }
    }
  });
}

// make a histogram!!!!!! obviously
function makeHistogram(nums, plotId, title) {
  //positions and dimensions
  var margin = {
    top: 100,
    right: 200,
    bottom: 30,
    left: 130
  };
  var width = 700;
  var height = 500;

  updateDimensions(window.innerWidth, margin);

  // Add the svg canvas
  var svg = d3.select("#" + plotId)
      .append("svg")
          .attr("width", width + margin.left + margin.right)
          .attr("height", height + margin.top + margin.bottom)
      .attr("id", "histogram" + plotId)
      .append("g")
          .attr("transform",
                "translate(" + margin.left + "," + margin.top + ")")
      .attr("fill", "white");

  // for resizing
  d3.select("#histogram" + plotId)
  .attr( 'preserveAspectRatio',"xMinYMin meet")
  .attr("viewBox", "0 0 800 400")
  .attr('width', '700')

  // graph title
  svg.append("text")
    .attr("text-anchor", "middle")
    .attr("fill","black")
    .attr("dy", 0)
    .style("font-size", "23px")
    .attr("transform", "translate("+ (width/2.2) +","+ -50 +")")
    .text(title)
    .call(wrap, width)

  // set the ranges
  var xrange = d3.scale.linear()
    .domain([1, d3.max(nums)])
    .range([0, width]);

  //format x axis
  var xAxis = d3.svg.axis()
    .scale(xrange)
    .orient("bottom");

  //make x axis
  svg.append("g")
      .attr("class", "x axis")
      .attr("fill", "none")
      .attr("stroke", "black")
      .attr("stroke-width", 1)
      .attr("transform", "translate(0," + height + ")")
      .call(xAxis);

  // X axis label
  svg.append("text")
      .attr("text-anchor", "middle")
      .attr("fill","black")
      .attr("transform", "translate("+ (width/2) +","+(height+50)+")")
      .text("# targets");

  //params for histogram
  var bins = d3.layout.histogram()
    .bins(xrange.ticks(25))
    (nums);


  var yrange = d3.scale.linear()
    .domain([0,d3.max(bins, function(d) {return d.length })])
    .range([height, 0]);

  //format y axis
  var yAxis = d3.svg.axis()
    .scale(yrange)
    .orient("left")
    .tickSize(5)
    .tickFormat(d3.format("s"));

  // Add the Y Axis
  svg.append("g")
      .attr("class", "y axis")
      .attr("fill", "none")
      .attr("stroke", "black")
      .attr("stroke-width", 1)
      .call(yAxis)

  // Y axis label
  svg.append("text")
      .attr("text-anchor", "middle")
      .attr("fill","black")
      .attr("transform", "translate("+ -60 +","+(height/2)+")rotate(-90)")
      .text("# compounds");

  svg.selectAll(".tick")
    .select("text")
    .attr("fill", "black")
    .attr("stroke", "none")


  // group for bars
  var bar = svg.selectAll(".bar")
    .data(bins)
    .enter().append("g")
      .attr("class", "bar")
      .attr("transform", function(d) { return "translate(" + xrange(d.x) + "," + yrange(d.y) + ")"; });

  bar.append("rect")
      .attr("x", 1)
      .attr("width", (xrange(bins[0].dx) - xrange(0)) -1)
      .attr("height", function(d) { return height  - yrange(d.y); })
      .attr("fill", "#207cc1")
      .on({
        "mouseover": function(d,i) {
          d3.select(this).transition()
          .duration(300).style("opacity", 0.6);
          d3.select("#" + "label" + i).transition()
          .duration(300).style("opacity", 1);
          d3.select(this).style("cursor", "pointer");
        },
        "mouseout": function(d,i) {
          d3.select(this).transition()
          .duration(300).style("opacity", 1);
          d3.select("#" + "label" + i).transition()
          .duration(300).style("opacity", 0);
          d3.select(this).style("cursor", "default");
        }
      })

  // bar tooltip
  bar.append("text")
      .attr("dy", ".75em")
      .attr("y", -12)
      .attr("fill", "black")
      .attr("id", function(d,i) {return "label" + i})
      .attr("font-size", "12px")
      .style("opacity", 0)
      .attr("x", (xrange(bins[0].dx) - xrange(0)) / 2)
      .attr("text-anchor", "middle")
      .text(function(d) { return d.y; });


      var arrTemp = title.split(" ");
      var descriptor = arrTemp[0] + "-" + arrTemp[1] + "-" + arrTemp[2] + "-" + arrTemp[3]

      d3.select("#" + plotId).append("button")
          .attr("type","button")
          .attr("id", "button" + plotId)
          .attr("class", "downloadButton")
          .text("Download SVG")
          .on("click", function() {
              // download the svg
              downloadSVG(plotId, descriptor);
          });

}
