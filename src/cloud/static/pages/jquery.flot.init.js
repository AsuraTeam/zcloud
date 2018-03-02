/**
 * Theme: Adminox Admin Template
 * Author: Coderthemes
 * Module/App: Flot-Chart
 */

! function($) {
	"use strict";

	var FlotChart = function() {
		this.$body = $("body")
		this.$realData = []
	};

	//creates plot graph
	FlotChart.prototype.createPlotGraph = function(selector, data1, data2, data3, labels, colors, borderColor, bgColor) {
		//shows tooltip
		function showTooltip(x, y, contents) {
			$('<div id="tooltip" class="tooltipflot">' + contents + '</div>').css({
				position : 'absolute',
				top : y + 5,
				left : x + 5
			}).appendTo("body").fadeIn(200);
		}


		$.plot($(selector), [{
			data : data1,
			label : labels[0],
			color : colors[0]
		}, {
			data : data2,
			label : labels[1],
			color : colors[1]
		},
			{
			data : data3,
			label : labels[2],
			color : colors[2]
		}], {
			series : {
				lines : {
					show : true,
					fill : true,
					lineWidth : 2,
					fillColor : {
						colors : [{
							opacity : 0.3
						}, {
							opacity : 0.3
						}, {
							opacity: 0.3
						}]
					}
				},
				points : {
					show : true
				},
				shadowSize : 0
			},

			grid : {
				hoverable : true,
				clickable : true,
				borderColor : borderColor,
				tickColor : "#f9f9f9",
				borderWidth : 1,
				labelMargin : 30,
				backgroundColor : bgColor
			},
			legend : {
				position: "ne",
				margin : [0, -32],
				noColumns : 0,
				labelBoxBorderColor : null,
				labelFormatter : function(label, series) {
					// just add some space to labes
					return '' + label + '&nbsp;&nbsp;';
				},
				width : 30,
				height : 2
			},
			yaxis : {
				axisLabel: "Daily Visits",
				tickColor : '#f5f5f5',
				font : {
					color : '#bdbdbd'
				}
			},
			xaxis : {
				axisLabel: "Last Days",
				tickColor : '#f5f5f5',
				font : {
					color : '#bdbdbd'
				}
			},
			tooltip : true,
			tooltipOpts : {
				content : '%s: Value of %x is %y',
				shifts : {
					x : -60,
					y : 25
				},
				defaultTheme : false
			},
			splines: {
				show: true,
				tension: 0.1, // float between 0 and 1, defaults to 0.5
				lineWidth: 1 // number, defaults to 2
			}
		});
	},
	//end plot graph

	//creates plot Dot graph
	FlotChart.prototype.createPlotDotGraph = function(selector, data1, data2, labelsDot, colorsDot, borderColorDot, bgColorDot) {
		//shows tooltip
		function showTooltip(x, y, contents) {
			$('<div id="tooltip" class="tooltipflot">' + contents + '</div>').css({
				position : 'absolute',
				top : y + 5,
				left : x + 5
			}).appendTo("body").fadeIn(200);
		}


		$.plot($(selector), [{
			data : data1,
			label : labelsDot[0],
			color : colorsDot[0]
		}, {
			data : data2,
			label : labelsDot[1],
			color : colorsDot[1]
		}],
			{
			series : {
				lines : {
					show : true,
					fill : false,
					lineWidth : 3,
					fillColor : {
						colors : [{
							opacity : 0.3
						}, {
							opacity : 0.3
						}]
					}
				},
				curvedLines: {
					apply: true,
					active: true,
					monotonicFit: true
				},
				splines: {
					show: true,
					tension: 0.4,
					lineWidth: 5,
					fill: 0.4
				}
			},

			grid : {
				hoverable : true,
				clickable : true,
				borderColor : borderColorDot,
				tickColor : "#f9f9f9",
				borderWidth : 1,
				labelMargin : 10,
				backgroundColor : bgColorDot
			},
			legend : {
				position : "ne",
				margin : [0, -32],
				noColumns : 0,
				labelBoxBorderColor : null,
				labelFormatter : function(label, series) {
					// just add some space to labes
					return '' + label + '&nbsp;&nbsp;';
				},
				width : 30,
				height : 2
			},
			yaxis : {
				axisLabel: "Gold Price(USD)",
				tickColor : '#f5f5f5',
				font : {
					color : '#bdbdbd'
				}
			},
			xaxis : {
				axisLabel: "Numbers",
				tickColor : '#f5f5f5',
				font : {
					color : '#bdbdbd'
				}
			},
			tooltip : false,
		});
	},
	//end plot Dot graph

	//creates Pie Chart
	FlotChart.prototype.createPieGraph = function(selector, labels, datas, colors) {
		var data = [{
			label : labels[0],
			data : datas[0]
		}, {
			label : labels[1],
			data : datas[1]
		}, {
			label : labels[2],
			data : datas[2]
		},{
			label : labels[3],
			data : datas[3]
		}];
		var options = {
			series : {
				pie: {
					show: true,
					radius: 1,
					label: {
						show: true,
						radius: 1,
						background: {
							opacity: 0.2
						}
					}
				}
			},
			legend : {
				show : false
			},
			grid : {
				hoverable : true,
				clickable : true
			},
			colors : colors,
			tooltip : true,
			tooltipOpts : {
				content : "%s, %p.0%"
			}
		};

		$.plot($(selector), data, options);
	},

	//returns some random data
	FlotChart.prototype.randomData = function() {
		var totalPoints = 300;
		if (this.$realData.length > 0)
			this.$realData = this.$realData.slice(1);

		// Do a random walk
		while (this.$realData.length < totalPoints) {

			var prev = this.$realData.length > 0 ? this.$realData[this.$realData.length - 1] : 50,
			    y = prev + Math.random() * 10 - 5;

			if (y < 0) {
				y = 0;
			} else if (y > 100) {
				y = 100;
			}

			this.$realData.push(y);
		}

		// Zip the generated y values with the x values
		var res = [];
		for (var i = 0; i < this.$realData.length; ++i) {
			res.push([i, this.$realData[i]])
		}

		return res;
	}, FlotChart.prototype.createRealTimeGraph = function(selector, data, colors) {
		var plot = $.plot(selector, [data], {
			colors : colors,
			series : {
				grow : {
					active : false
				}, //disable auto grow
				shadowSize : 0, // drawing is faster without shadows
				lines : {
					show : true,
					fill : true,
					lineWidth : 2,
					steps : false
				}
			},
			grid : {
				show : true,
				aboveData : false,
				color : '#dcdcdc',
				labelMargin : 15,
				axisMargin : 0,
				borderWidth : 0,
				borderColor : null,
				minBorderMargin : 5,
				clickable : true,
				hoverable : true,
				autoHighlight : false,
				mouseActiveRadius : 20
			},
			tooltip : true, //activate tooltip
			tooltipOpts : {
				content : "Value is : %y.0" + "%",
				shifts : {
					x : -30,
					y : -50
				}
			},
			yaxis : {
				axisLabel: "Response Time (ms)",
				min : 0,
				max : 100,
				tickColor : '#f5f5f5',
				color : 'rgba(0,0,0,0.1)'
			},
			xaxis : {
				axisLabel: "Point Value (1000)",
				show : true,
				tickColor : '#f5f5f5'
			}
		});

		return plot;
	},
	//creates Donut Chart
	FlotChart.prototype.createDonutGraph = function(selector, labels, datas, colors) {
		var data = [{
			label : labels[0],
			data : datas[0]
		}, {
			label : labels[1],
			data : datas[1]
		}, {
			label : labels[2],
			data : datas[2]
		}, {
			label : labels[3],
			data : datas[3]
		}];
		var options = {
			series : {
				pie : {
					show : true,
					innerRadius : 0.7
				}
			},
			legend : {
				show : true,
				labelFormatter : function(label, series) {
					return '<div style="font-size:14px;">&nbsp;' + label + '</div>'
				},
				labelBoxBorderColor : null,
				margin : 50,
				width : 20
			},
			grid : {
				hoverable : true,
				clickable : true
			},
			colors : colors,
			tooltip : true,
			tooltipOpts : {
				content : "%s, %p.0%"
			}
		};

		$.plot($(selector), data, options);
	},
	//creates Bar Chart
	FlotChart.prototype.createStackBarGraph = function(selector, ticks, colors, data) {
		var options = {
			bars: {
				show: true,
				barWidth: 0.2,
				fill: 1
			},
			grid: {
				show: true,
				aboveData: false,
				labelMargin: 5,
				axisMargin: 0,
				borderWidth: 1,
				minBorderMargin: 5,
				clickable: true,
				hoverable: true,
				autoHighlight: false,
				mouseActiveRadius: 20,
				borderColor: '#f5f5f5'
			},
			series: {
				stack: 0
			},
			legend: {
				position: "ne",
				margin: [0, -32],
				noColumns: 0,
				labelBoxBorderColor: null,
				labelFormatter: function (label, series) {
					// just add some space to labes
					return '' + label + '&nbsp;&nbsp;';
				},
				width: 30,
				height: 2
			},
			yaxis: ticks.y,
			xaxis: ticks.x,
			colors: colors,
			tooltip: true, //activate tooltip
			tooltipOpts: {
				content: "%s : %y.0",
				shifts: {
					x: -30,
					y: -50
				}
			}
		};
		$.plot($(selector), data, options);
	},
	//creates Line Chart
	FlotChart.prototype.createLineGraph = function(selector, ticks, colors, data) {
		var options = {
			series: {
				lines: {
					show: true
				},
				points: {
					show: true
				}
			},
			legend : {
				position : "ne",
				margin : [0, -32],
				noColumns : 0,
				labelBoxBorderColor : null,
				labelFormatter : function(label, series) {
					// just add some space to labes
					return '' + label + '&nbsp;&nbsp;';
				},
				width : 30,
				height : 2
			},
			yaxis: ticks.y,
			xaxis: ticks.x,
			colors: colors,
			grid: {
				hoverable: true,
				borderColor: '#f5f5f5',
				borderWidth: 1,
				backgroundColor: '#fff'
			},
			tooltip: true, //activate tooltip
			tooltipOpts: {
				content: "%s : %y.0",
				shifts: {
					x: -30,
					y: -50
				}
			}
		};
		return $.plot($(selector), data, options);
	},
	//creates Combine Chart
	FlotChart.prototype.createCombineGraph = function(selector, ticks, labels, datas) {

		var data = [{
			label : labels[0],
			data : datas[0],
			lines : {
				show : true,
				fill : true
			},
			points : {
				show : true
			}
		}, {
			label : labels[1],
			data : datas[1],
			lines : {
				show : true
			},
			points : {
				show : true
			}
		}, {
			label : labels[2],
			data : datas[2],
			bars : {
				show : true
			}
		}];
		var options = {
			series : {
				shadowSize : 0
			},
			grid : {
				hoverable : true,
				clickable : true,
				tickColor : "#f9f9f9",
				borderWidth : 1,
				borderColor : "#eeeeee"
			},
			colors : ['#4489e4', '#78c350', '#f18f1c'],
			tooltip : true,
			tooltipOpts : {
				defaultTheme : false
			},
			legend : {
				position : "ne",
				margin : [0, -32],
				noColumns : 0,
				labelBoxBorderColor : null,
				labelFormatter : function(label, series) {
					// just add some space to labes
					return '' + label + '&nbsp;&nbsp;';
				},
				width : 30,
				height : 2
			},
			yaxis : {
				axisLabel: "Point Value (1000)",
				tickColor : '#f5f5f5',
				font : {
					color : '#bdbdbd'
				}
			},
			xaxis : {
				axisLabel: "Daily Hours",
				ticks: ticks,
				tickColor : '#f5f5f5',
				font : {
					color : '#bdbdbd'
				}
			}
		};

		$.plot($(selector), data, options);
	},

	//initializing various charts and components
	FlotChart.prototype.init = function() {
		//plot graph data
		var uploads = [[0, 13], [1, 13], [2, 14], [3, 17], [4, 13], [5, 10], [6, 12],[7, 13], [8, 12], [9, 20], [10, 18], [11, 16], [12, 14]];
		var downloads = [[0, 8], [1, 10], [2, 12], [3, 14], [4, 11], [5, 7], [6, 9],[7, 10], [8, 9], [9, 17], [10, 15], [11, 13], [12, 11]];
		var downloads1 = [[0, 3], [1, 6], [2, 8], [3, 10], [4, 7], [5, 3], [6, 5],[7, 7], [8, 6], [9, 14], [10, 12], [11, 10], [12, 8]];
		var plabels = ["Google", "Yahoo", "Facebook"];
		var pcolors = ['#4489e4', '#78c350', '#f18f1c'];
		var borderColor = '#f5f5f5';
		var bgColor = '#fff';
		this.createPlotGraph("#website-stats", uploads, downloads,downloads1, plabels, pcolors, borderColor, bgColor);

		//plot graph Dot data
		var uploadsDot = [[0, 2], [1, 4], [2, 7], [3, 9], [4, 6], [5, 3], [6, 10],[7, 8], [8, 5], [9, 14], [10, 10], [11, 10], [12, 8]];
		var downloadsDot = [[0, 1], [1, 3], [2, 6], [3, 7], [4, 4], [5, 2], [6, 8],[7, 6], [8, 4], [9, 10], [10, 8], [11, 14], [12, 5]];
		var plabelsDot = ["Visits", "Page views"];
		var pcolorsDot = ['#34bd25','#ff4433'];
		var borderColorDot = '#f5f5f5';
		var bgColorDot = '#fff';
		this.createPlotDotGraph("#website-stats1", uploadsDot, downloadsDot, plabelsDot, pcolorsDot, borderColorDot, bgColorDot);

		//Pie graph data
		var pielabels = ["Mobile Phones", "iPhone 7", "iMac", "Macbook"];
		var datas = [20, 30, 15, 32];
		var colors = ['#47dea3', '#23d0ea','#fb6d9d','#3f4958'];
		this.createPieGraph("#pie-chart #pie-chart-container", pielabels, datas, colors);

		//real time data representation
		var plot = this.createRealTimeGraph('#flotRealTime', this.randomData(), ['#64c5b1']);
		plot.draw();
		var $this = this;
		function updatePlot() {
			plot.setData([$this.randomData()]);
			// Since the axes don't change, we don't need to call plot.setupGrid()
			plot.draw();
			setTimeout(updatePlot, $('html').hasClass('mobile-device') ? 500 : 500);
		}

		updatePlot();

		//Donut pie graph data
		var donutlabels = ["Mobile Phones", "iPhone 7", "iMac", "Macbook"];
		var donutdatas = [35, 20, 10, 20];
		var donutcolors = ['#60befc', '#6248ff','#e3b0db','#dbdbdb'];
		this.createDonutGraph("#donut-chart #donut-chart-container", donutlabels, donutdatas, donutcolors);

		//Combine graph data
		var data24Hours = [[0, 201], [1, 520], [2, 337], [3, 261], [4, 157], [5, 95], [6, 200], [7, 250], [8, 320], [9, 500], [10, 152], [11, 214], [12, 364], [13, 449], [14, 558], [15, 282], [16, 379], [17, 429], [18, 518], [19, 470], [20, 330], [21, 245], [22, 358], [23, 74]];
		var data48Hours = [[0, 311], [1, 630], [2, 447], [3, 371], [4, 267], [5, 205], [6, 310], [7, 360], [8, 430], [9, 610], [10, 262], [11, 324], [12, 474], [13, 559], [14, 668], [15, 392], [16, 489], [17, 539], [18, 628], [19, 580], [20, 440], [21, 355], [22, 468], [23, 184]];
		var dataDifference = [[23, 727], [22, 128], [21, 110], [20, 92], [19, 172], [18, 63], [17, 150], [16, 592], [15, 12], [14, 246], [13, 52], [12, 149], [11, 123], [10, 2], [9, 325], [8, 10], [7, 15], [6, 89], [5, 65], [4, 77], [3, 600], [2, 200], [1, 385], [0, 200]];
		var ticks = [[0, "22h"], [1, ""], [2, "00h"], [3, ""], [4, "02h"], [5, ""], [6, "04h"], [7, ""], [8, "06h"], [9, ""], [10, "08h"], [11, ""], [12, "10h"], [13, ""], [14, "12h"], [15, ""], [16, "14h"], [17, ""], [18, "16h"], [19, ""], [20, "18h"], [21, ""], [22, "20h"], [23, ""]];
		var combinelabels = ["Last 24 Hours", "Last 48 Hours", "Difference"];
		var combinedatas = [data24Hours, data48Hours, dataDifference];

		this.createCombineGraph("#combine-chart #combine-chart-container", ticks, combinelabels, combinedatas);

		//bar chart = stacked
		var stack_ticks = {
			y: {
				axisLabel: "Sales Value (USD)",
				tickColor: '#f5f5f5',
				font: {
					color: '#bdbdbd'
				}
			},
			x: {
				axisLabel: "Last 10 Days",
				tickColor: '#f5f5f5',
				font: {
					color: '#bdbdbd'
				}
			}
		};

		//random data
		var d1 = [];
		for (var i = 0; i <= 10; i += 1)
			d1.push([i, parseInt(Math.random() * 30)]);

		var d2 = [];
		for (var i = 0; i <= 10; i += 1)
			d2.push([i, parseInt(Math.random() * 30)]);

		var d3 = [];
		for (var i = 0; i <= 10; i += 1)
			d3.push([i, parseInt(Math.random() * 30)]);

		var ds = new Array();

		ds.push({
			label: "Series One",
			data: d1,
			bars: {
				order: 3
			}
		});
		ds.push({
			label: "Series Two",
			data: d2,
			bars: {
				order: 2
			}
		});
		ds.push({
			label: "Series Three",
			data: d3,
			bars: {
				order: 1
			}
		});
		this.createStackBarGraph("#ordered-bars-chart", stack_ticks, ['#23d0ea', '#3f4958', "#ebeff2"], ds);


		//creating line chart
		var line_ticks = {
			y: {
				min: -1.2,
				max: 1.2,
				tickColor: '#f5f5f5',
				font : {
					color : '#bdbdbd'
				}
			},
			x: {
				tickColor: '#f5f5f5',
				font : {
					color : '#bdbdbd'
				}
			}
		};

		//sample data
		var sin = [],
			cos = [];
		var offset = 0;
		for (var i = 0; i < 12; i += 0.2) {
			sin.push([i, Math.sin(i + offset)]);
			cos.push([i, Math.cos(i + offset)]);
		}
		var line_data = [
			{
				data: sin,
				label: "Google",
			},
			{
				data: cos,
				label: "Yahoo"
			}
		];
		this.createLineGraph("#line-chart-alt", line_ticks, ["#f5707a", "#188ae2"], line_data);
	},

	//init flotchart
	$.FlotChart = new FlotChart, $.FlotChart.Constructor =
	FlotChart

}(window.jQuery),

//initializing flotchart
function($) {
	"use strict";
	$.FlotChart.init()
}(window.jQuery);