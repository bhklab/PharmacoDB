//takes in an array of stddevs, and an array of means
//and returns an array of 95% confidence intervals
function confidInterval(stddev, mean) {
  return [mean + (1.96*stddev), mean - (1.96*stddev)]
}


//var p1 = [0.05659940, 0.08041887, 398.00000000, 0.48199413]
// var p2 = [ -0.07103750 ,  0.08601112 ,404.00000000 ,  0.40937711  ]
// var p3 = [1.330870e-01, 4.840363e-02, 8.470000e+02, 6.100783e-03 ]
// var params = [p1, p2, p3]
// //estimate, stderr, n, pvalue
// makeForestPlot(params, "body");

function makeForestPlot(params, plotId) {
  //estimate, stderr, n, pvalue
  //positions and dimensions

  var width = 500;
  var height = 500;


  var names = ["GDSC", "gCSI", "CCLE"]
  names.unshift(" ")
  
  // formatting the given data for easier parsing
  var estimates = []
  var n = []
  var p = []
  var stderr = []
  var confids = []
  for (var i = 0; i < params.length; i++) {
    estimates.push(params[i][0])
    n.push(params[i][2])
    p.push(params[i][3])
    stderr.push(params[i][2])
    confids.push(confidInterval(params[i][1], params[i][0]))
  }
  estimates.unshift(0)
  n.unshift(0)
  p.unshift("")
  stderr.unshift(0)
  confids.unshift([0,0])

  //make the x axis fit the circles with stderr bars
  var tempMax = [], tempMin = []
  for (var i = 0; i < confids.length; i++) {
    tempMax.push( ((confids[i][0])) )
    tempMin.push( (confids[i][1]) )
  }
  
  //getting maxes for easy parsing
  var maxNum = d3.max(tempMax)
  var minNum = d3.min(tempMin)

  //set range for data by domain, and scale by range
  var xrange = d3.scale.linear()
    .domain([minNum, maxNum])
    .range([0, width])
    .nice();

  var yrange = d3.scale.linear()
    .domain([0, names.length])
    .range([0, height]);

  //set axes for graph
  var xAxis = d3.svg.axis()
    .scale(xrange)
    .innerTickSize(-height)
    .outerTickSize(5)
    .orient("bottom")
    .tickPadding(10)

  var yAxis = d3.svg.axis()
    .scale(yrange)
    .orient("left")
    .innerTickSize(-width)
    .outerTickSize(5)
		.tickFormat(function(d,i){ return names[i] })
		.tickValues(d3.range(names.length));

  // Add the svg canvas
  var svg = d3.select(plotId) //TODO: "#" + plotId
      .append("svg")
        .attr("fill", "white")
        .attr("transform", "translate(100,50)")
          .attr("width", width + 300)
          .attr("height", height + 150)
      .attr("id", "forestPlot" + plotId)

  // graph title
  svg.append("text")
    .attr("text-anchor", "middle")
    .attr("fill","black")
    .attr("dx", 0)
    .attr("dy", 0)
    .style("font-size", "23px")
    .attr("transform", "translate("+ (width/2.2 + 100 ) +","+ 30 +")")
    .text("forest")

  // Add the X Axis
  var xAxis = svg.append("g")
      .attr("class", "x axis")
      .attr("transform", "translate(100," +  (height+50) + ")")
      .attr("fill", "none")
      .attr("stroke", "black")
      .attr("stroke-width", 1)
      .call(xAxis);

  d3.selectAll(".tick")
    .select("text")
    .attr("fill", "black")
    .attr("stroke", "none")


  // X axis label
  svg.append("text")
      .attr("text-anchor", "middle")
      .attr("fill","black")
      .attr("transform", "translate("+ (width/2 + 100) +","+(height+100)+")")
      .text("Regression Estimate");

  // Add the Y Axis
  var yAxis = svg.append("g")
      .attr("class", "y axis")
      .attr("fill", "none")
      .attr("stroke", "black")
      .attr("transform", "translate(100,50)")
      .attr("stroke-width", 1)
      .call(yAxis)
      .selectAll("text")
      .attr("fill", "blue")
      .attr("stroke", "none")
      .on("click", function(i){
            document.location.href = "/search?q=" + names[i]
        })
        .on({
          "mouseover": function(d) {
            d3.select(this).transition().duration(400).style("opacity", 0.5);
            d3.select(this).style("cursor", "pointer");
          },
          "mouseout": function(d) {
            d3.select(this).transition().duration(400).style("opacity", 1);
            d3.select(this).style("cursor", "default");
          }
        });

  
        d3.selectAll(".tick")
        .select("line")
        .style("opacity", 0.2)

    // Y axis label
    svg.append("text")
    .attr("text-anchor", "middle")
    .attr("fill","black")
    .attr("transform", "translate("+ 30 +","+(height/1.7)+")rotate(-90)")
    .text("Study");

    // adding circles with proportionate size
    var circles = svg.selectAll("circle")
      .data(estimates)
      .enter();

      circles.append("circle")
        .attr("id", function(d, i) {return "circle" + i})
        .attr("r", function(d,i) {return n[i]/30})
        .attr("fill", "black")
        .attr("cx", function(d,i) {return xrange(d)})
        .attr("cy", function(d,i) {return yrange(i)})
        .attr("transform", "translate(100,50)")

      //pvals
      circles.append("text")
        .attr("fill", "black")
        .attr("class", function(d,i) {return "pval" + i})
        .attr("dx", function(d,i) {return xrange(d) + 100})
        .attr("dy", function(d,i) {return yrange(i) + 10})
        .attr("text-anchor", "middle")
        .text(function(d,i) { return "P: " + p[i]})

      d3.select(".pval0").remove()

    // error bars
    var errorBar = svg.selectAll("line.error")
        .data(estimates)
        .enter();

        errorBar.append("line")
          .attr("class", "error")
          .style("stroke", "black")
          .attr("x1", function(d,i) { return xrange(confids[i][1]) })
          .attr("x2", function(d,i) { return xrange(confids[i][0]) })
          .attr('y1', function(d,i) { return yrange(i);})
          .attr('y2', function(d,i) { return yrange(i);})
          .attr("transform", "translate(100,50)");



    d3.select("#" + plotId).append("button")
        .attr("type","button")
        .attr("id", "button" + plotId)
        .attr("class", "downloadButton")
        .text("Download SVG")
        .on("click", function() {
            // download the svg
            if (cell_line != "") {
              downloadSVG(plotId, cell_line);
            }
            else {
              downloadSVG(plotId, xAxisLabel)
            }
        });
}
