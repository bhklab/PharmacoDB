function updateDimensions(winWidth, margin) {
  width = winWidth - margin.left - margin.right;
  height = 500 - margin.top - margin.bottom;
}

String.prototype.replaceAll = String.prototype.replaceAll || function(s, r) {
  return this.replace(new RegExp(s, 'g'), r);
};

function cleanNames(names) {
  var newNames = [];
  for (var i = 0; i < names.length; i++) {
    newNames.push(names[i].replaceAll("_", " "))
  }
  return newNames;
}

function makePie(names, nums, other_names, other_nums) {
  //setting up data: [{"name": "tissueName", "num": 3}]
  var cleanedNames = cleanNames(names);
  var data = [];
  for (var i = 0; i < names.length; i++) {
    data.push({"name": names[i], "num": nums[i]});
  }

  //sum of cell lines
  var sum = nums.reduce(function (a, b) {
    return a + b;
  }, 0);

  //positions and dimensions
  var margin = {
    top: 0,
    right: 1000,
    bottom: 120,
    left: 10
  };
  // var width = parseInt(d3.select('body').style('width'), 10);
  // var height = parseInt(d3.select('body').style('height'), 10);
  var width = 850;
  var height = 400;
  var radius = 450;
  var color = d3.scale.category20();

  // reverse the tissue names for the legend, not sure why I have to do this
  // but it must be done9
  var revArray = names.reverse();

  updateDimensions(window.innerWidth, margin);

  // Add the svg canvas
  var svg = d3.select("#count_stat")
      .append("svg:svg")
        .data([data])
          .attr("width", width + margin.left + margin.right)
          .attr("height", height + margin.top + margin.bottom)
          .attr("transform",
                "translate(" + 0 + "," + margin.top + ")")
      .attr("id", "pie")
      .append("g")
          .attr("transform",
                "translate(" + radius + "," + (radius+120) + ")")
      .attr("fill", "white");

  // for resizing
  d3.select("#pie")
  .attr( 'preserveAspectRatio',"xMinYMin meet")
  .attr("viewBox", "0 0 1420 400")
  .attr('width', '700')

  // graph title
  svg.append("text")
    .attr("text-anchor", "middle")
    .attr("fill","black")
    .style("font-size", "40px")
    .attr("transform", "translate("+ (width/4) +","+ -490 +")")
    .text("Relative Percentage of Cell lines Per Tissue in PharmacoDB");

  var arc = d3.svg.arc()
    .outerRadius(radius)
    .innerRadius(radius - 200);
      //.outerRadius(radius);

  var pie = d3.layout.pie()
      .value(function(d) {
          return d.num;
      });

  //making each slice of pie
  var arcs = svg.selectAll("g.slice")
      .data(pie)
      .enter()
          .append("svg:g")
          .attr("class", function(d,i) {return data[i].name});


      arcs.append("a")
        .attr("xlink:href", function(d,i) {
          if(data[i].name != "Other") {return "/search?q=" + data[i].name}
        })
        .append("svg:path")
        .attr("fill", function(d, i) { return color(i); } )
        .attr("d", arc);

      arcs.append("svg:title")
        .text(function(d,i) {return data[i].name});

      //for other section, show table
      var other = d3.select(".Other");
      var active;
      other.on("click", function() {
        active   = other.active ? false : true,
          newOpacity = active ? 1 : 0;
        d3.select("#otherTissues").transition().duration(500).style("opacity", newOpacity);

        if (active) {
          d3.select("#legendLabelOther").transition().duration(500).attr("fill", "silver");
          d3.select(".Other").transition().duration(500).style("opacity", 0.7);
        } else {
          d3.select("#legendLabelOther").transition().duration(500).attr("fill", "black");
          d3.select(".Other").transition().duration(500).style("opacity", 1);
        }
         other.active = active;
      })
      .on({
        "mouseover": function(d) {
          d3.select(".Other").transition()
          .duration(400).style("opacity", 0.7);
          d3.select("#legendLabelOther").transition().duration(500).attr("fill", "silver");
          d3.select(this).style("cursor", "pointer");
        },
        "mouseout": function(d) {
          d3.select(".Other").transition()
          .duration(400).style("opacity", 1);
          d3.select("#legendLabelOther").transition().duration(500).attr("fill", "black");
          d3.select(this).style("cursor", "default");
        }
      });

      for (var i = 0; i < names.length; i++) {
        if (data[i].name != "Other") {
          d3.select("." + data[i].name).on({
            "mouseover": function() {
              d3.select(this).transition()
              .duration(400).style("opacity", 0.7); // slice
              d3.select(this).style("cursor", "pointer");
            },
            "mouseout": function() {
              d3.select(this).transition()
              .duration(400).style("opacity", 1); // slice opacity
              d3.select(this).style("cursor", "default");
            }
          })
        }
      }

      //percentage
      arcs.append("svg:text")
          .attr("transform", function(d) {  //center
          d.innerRadius = 0;
          d.outerRadius = radius;
          var coords = arc.centroid(d);
          coords[0] *= 1;
          coords[1] *= 1;
          return "translate(" + coords + ")"; //return coordinates
        })
        //.attr("dy", "1em") for inside
      .attr("text-anchor", "middle")
      .attr("fill", "white")
      .style("font-size", "25px")
      .text(function(d, i) {
        return d3.format(".2f")((data[i].num/sum)*100) + "%";
      });

      d3.select("#count_stat").append("button")
          .attr("type","button")
          .attr("id", "buttoncount_stat")
          .attr("class", "downloadButton")
          .text("Download SVG")
          .on("click", function() {
              // download the svg
              downloadSVG("count_stat", "tissue");
          });

      // make a table that for the "Other" tissues
      // make sure there is enough data to populate each subarray of 3
      var tableData = [];
      var addElems = other_names.length % 3;
      other_names_copy = other_names.slice(0);
      other_nums_copy = other_nums.slice(0);
      for (var i = 0; i < 3-addElems; i++) {
        other_names_copy.push(" ");
        other_nums_copy.push(0);
      }

      //make multidimensional array of 3
      for (var i = 0; i < other_names.length; i += 3) {
        var temp = []
        for (var j = 0; j < 3; j++) {
          temp.push({"name": other_names_copy[i+j], "num": other_nums_copy[i+j]})
          //temp.push(other_names_copy[i+j])
        }
        tableData.push(temp)
      }

      var table = d3.select("#count_stat").append('table')
        .attr("id", "otherTissues")
        .attr("class", "table table-bordered")
        .style("border", "1px solid silver")
        .style("width", "1000px")
        .style("opacity", "0")
        .style("margin", "30px 0px 0px -130px")

      var tbody = table.append("tbody");

      //create table rows
      var tr = tbody.selectAll("tr")
        .data(tableData)
        .enter()
        .append("tr")
        //.style("border", "1px solid silver");

      // create table cells
      var td = tr.selectAll("td")
        .data(function(d) {return d;})
        .enter()
        .append("td")
        //.style("font-size", "12px")
        //.style("border", "1px solid silver")
        .attr("width", 400)
        .style("padding-top", "5px")
        .style("padding-bottom", "5px")
        .html(function(d,i){
          //extremely tedious parsing for data
          if ((d.name != "" && d.num != 0) || (d.num == 0 && d.name != " ") || (d.num != 0)) {
            return "<a href=\"" + "/search?q=" + d.name + "\">" + d.name.replaceAll("_", " ") + "</a>" + " (" + d3.format(".2f")((d.num/sum)*100) + "%)"
          }
        });


        var newCleanedNames = cleanedNames.slice(0);
        newCleanedNames.reverse();

    //legend
    for (var i = 0; i < data.length; i++) {
      svg.append('rect')
          .attr("x", width-350)
          .attr("y", 160 - i * 35)
          .attr("width", 15)
          .attr("height", 15)
          .style("fill", color(i));



      svg.append('text')
          .attr("x", width-330)
          .attr("y",  i * 35 - 245)
          .attr("id", "legendLabel" + revArray[i])
          .style("text-anchor", "start")
          .style("font-size", "25px")
          .style("opacity", 1)
          .attr("fill", "black")
          .on({
            "mouseover": function(d) {
              d3.select(this).transition()
              .duration(300).style("opacity", 0.6);
              d3.select(this).style("cursor", "pointer");
            },
            "mouseout": function(d) {
              d3.select(this).transition()
              .duration(300).style("opacity", 1);
              d3.select(this).style("cursor", "default");
            }
          })
          .html(function(d){
            if (names[i] == "Other") {
              return "Other (click for more)"
            } else {
              return "<a href=\"" + "/search?q=" + names[i] + "\">" + newCleanedNames[i] + "</a>";

            }
          });

          var other = d3.select("#legendLabelOther");
          var active;
          other.on("click", function() {
            active = other.active ? false : true,
              newOpacity = active ? 1 : 0;
            d3.select("#otherTissues").transition().duration(500).style("opacity", newOpacity);

            if (active) {
              d3.select("#legendLabelOther").transition().duration(500).attr("fill", "silver");
              d3.select(".Other").transition().duration(500).style("opacity", 0.7);
            } else {
              d3.select("#legendLabelOther").transition().duration(500).attr("fill", "black");
              d3.select(".Other").transition().duration(500).style("opacity", 1);
            }
             other.active = active;
          });
      }

      // total number of cell lines
      svg.append('text')
          .attr("x", width)
          .attr("y",  height + 40)
          .style("text-anchor", "end")
          .style("font-size", "25px")
          .attr("fill", "black")
          .text("Total number of cell lines: " + sum)



}
