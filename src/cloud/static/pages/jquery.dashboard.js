/**
* Theme: Adminox Dashboard
* Author: Coderthemes
* Dashboard
*/

jQuery(function($) {

  'use strict';

  var AdminoxAdmin = window.AdminoxAdmin || {};




  /*--------------------------------
   Window Based Layout
   --------------------------------*/
  AdminoxAdmin.dashboardEcharts = function() {


    /*--------------- Chart 1 -------------*/
    if($("#platform_type_dates_donut").length){
      var myChart = echarts.init(document.getElementById('platform_type_dates_donut'));

      var idx = 1;
      var option_dt = {

        timeline : {
          show: true,
          data : ['06-16','05-16','04-16'],
          label : {
            formatter : function(s) {
              return s.slice(0, 5);
            }
          },
          x:10,
          y:null,
          x2:10,
          y2:0,
          width:250,
          height:50,
          backgroundColor:"rgba(0,0,0,0)",
          borderColor:"#eaeaea",
          borderWidth:0,
          padding:5,
          controlPosition:"left",
          autoPlay:true,
          loop:true,
          playInterval:2000,
          lineStyle:{
            width:1,
            color:"#bdbdbd",
            type:""
          },

        },

        options : [
          {
            color: ['#dddddd','#64c5b1','#414b4f','#ee4b82','#45bbe0'],
            title : {
              text: '',
              subtext: ''
            },
            tooltip : {
              trigger: 'item',
              formatter: "{a} <br/>{b} : {c} ({d}%)"
            },
            legend: {
              show: false,
              x: 'left',
              orient:'vertical',
              padding: 0,
              data:['iPhone 7','Windows','Desktop','Mobiles','Others']
            },
            toolbox: {
              show : true,
              color : ['#bdbdbd','#bdbdbd','#bdbdbd','#bdbdbd'],
              feature : {
                mark : {show: false},
                dataView : {show: false, readOnly: true},
                magicType : {
                  show: true,
                  type: ['pie', 'funnel'],
                  option: {
                    funnel: {
                      x: '10%',
                      width: '80%',
                      funnelAlign: 'center',
                      max: 50
                    },
                    pie: {
                      roseType : 'none',
                    }
                  }
                },
                restore : {show: false},
                saveAsImage : {show: true}
              }
            },


            series : [
              {
                name:'06-16',
                type:'pie',
                radius : [20, '80%'],
                roseType : 'none',
                center: ['50%', '45%'],
                width: '50%',       // for funnel
                itemStyle : {
                  normal : { label : { show : true }, labelLine : { show : true } },
                  emphasis : { label : { show : false }, labelLine : {show : false} }
                },
                data:[{value: 35,  name:'iPhone 7'}, {value: 16,  name:'Windows'}, {value: 27,  name:'Desktop'}, {value: 29,  name:'Mobiles'}, {value: 12,  name:'Others'}]
              }
            ]
          },
          {
            series : [
              {
                name:'05-16',
                type:'pie',
                data:[{value: 42,  name:'iPhone 7'}, {value: 51,  name:'Windows'}, {value: 39,  name:'Desktop'}, {value: 25,  name:'Mobiles'}, {value: 9,  name:'Others'}]
              }
            ]
          },
          {
            series : [
              {
                name:'04-16',
                type:'pie',
                data:[{value: 29,  name:'iPhone 7'}, {value: 16,  name:'Windows'}, {value: 24,  name:'Desktop'}, {value: 19,  name:'Mobiles'}, {value: 5,  name:'Others'}]
              }
            ]
          },

        ] // end options object
      };

      myChart.setOption(option_dt);


    }


    /*-------------- Chart 2 ---------------*/
    if($("#user_type_bar").length){
// Initialize after dom ready
      var myChart = echarts.init(document.getElementById('user_type_bar'));

      var option = {

        // Setup grid
        grid: {
          zlevel: 0,
          x: 50,
          x2: 50,
          y: 20,
          y2: 20,
          borderWidth: 0,
          backgroundColor: 'rgba(0,0,0,0)',
          borderColor: 'rgba(0,0,0,0)',
        },

        // Add tooltip
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'shadow', // line|shadow
            lineStyle:{color: 'rgba(0,0,0,.5)', width: 1},
            shadowStyle:{color: 'rgba(0,0,0,.1)'}
          }
        },

        // Add legend
        legend: {
          data: []
        },
        toolbox: {
          orient: 'vertical',
          show : true,
          showTitle: true,
          color : ['#bdbdbd','#bdbdbd','#bdbdbd','#bdbdbd'],
          feature : {
            mark : {show: false},
            dataZoom : {
              show : true,
              title : {
                dataZoom : 'Data Zoom',
                dataZoomReset : 'Reset Zoom'
              }
            },
            dataView : {show: false, readOnly: true},
            magicType : {
              show: true,
              title : {
                bar : 'Bar',
                line : 'Area',
                stack : 'Stacked Bar',
                tiled: 'Tiled Bar'
              },
              type: ['bar','line','stack','tiled']
            },
            restore : {show: false},
            saveAsImage : {show: true,title:'Save as Image'}
          }
        },

        // Enable drag recalculate
        calculable: true,

        // Horizontal axis
        xAxis: [{
          type: 'category',
          boundaryGap: false,
          data: ['2016-06-01','2016-05-01','2016-04-01','2016-03-01','2016-02-01','2016-01-01','2015-12-01','2015-11-01','2015-10-01','2015-09-01'],
          axisLine: {
            show: true,
            onZero: true,
            lineStyle: {
              color: '#64c5b1',
              type: 'solid',
              width: '2',
              shadowColor: 'rgba(0,0,0,0)',
              shadowBlur: 5,
              shadowOffsetX: 3,
              shadowOffsetY: 3,
            },
          },
          axisTick: {
            show: false,
          },
          splitLine: {
            show: false,
            lineStyle: {
              color: '#fff',
              type: 'solid',
              width: 0,
              shadowColor: 'rgba(0,0,0,0)',
            },
          },
        }],

        // Vertical axis
        yAxis: [{
          type: 'value',
          splitLine: {
            show: false,
            lineStyle: {
              color: 'fff',
              type: 'solid',
              width: 0,
              shadowColor: 'rgba(0,0,0,0)',
            },
          },
          axisLabel: {
            show: false,
          },
          axisTick: {
            show: false,
          },
          axisLine: {
            show: false,
            onZero: true,
            lineStyle: {
              color: '#dddddd',
              type: 'solid',
              width: '0',
              shadowColor: 'rgba(0,0,0,0)',
              shadowBlur: 5,
              shadowOffsetX: 3,
              shadowOffsetY: 3,
            },
          },


        }],

        // Add series
        series: [
          {
            name: 'Registered Users',
            type: 'bar',
            smooth: true,
            symbol:'none',
            symbolSize:2,
            showAllSymbol: true,
            barWidth:10,
            barGap:'10%',
            itemStyle: {
              normal: {
                color:'#64c5b1',
                borderWidth:2, borderColor:'#64c5b1',
                areaStyle: {color:'#64c5b1', type: 'default'}
              }
            },

            data: [2323,2144,4534,1989,3232,2323,2144,4534,1989,3232,2323,2144,4534,1989,3232,2323,2144,4534,1989,3232,2323,2144,4534,1989,3232,2323,2144,4534,1989,3232]
          },
          {
            name: 'Guest Visitors',
            type: 'bar',
            smooth: true,
            symbol:'none',
            symbolSize:2,
            showAllSymbol: true,
            barWidth:10,
            barGap:'10%',
            itemStyle: {
              normal: {
                color:'#dddddd',
                borderWidth:2, borderColor:'#dddddd',
                areaStyle: {color:'#dddddd', type: 'default'}
              }
            },

            data: [5656,6567,7675,3423,4343,5656,6567,7675,3423,4343,5656,6567,7675,3423,4343,5656,6567,7675,3423,4343,5656,6567,7675,3423,4343,5656,6567,7675,3423,4343]
          },
        ]
      };

      // Load data into the ECharts instance
      myChart.setOption(option);

    }




    /*----------------- Chart 4 ------------------*/
    if($("#page_views_today").length){

// Initialize after dom ready
      var myChart = echarts.init(document.getElementById('page_views_today'));

      var option = {

        // Setup grid
        grid: {
          zlevel: 0,
          x: 20,
          x2: 20,
          y: 20,
          y2: 20,
          borderWidth: 0,
          backgroundColor: 'rgba(0,0,0,0)',
          borderColor: 'rgba(0,0,0,0)',
        },

        // Add tooltip
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'shadow', // line|shadow
            lineStyle:{color: 'rgba(0,0,0,.5)', width: 1},
            shadowStyle:{color: 'rgba(0,0,0,.1)'}
          }
        },

        // Add legend
        legend: {
          data: []
        },
        toolbox: {
          orient: 'vertical',
          show : true,
          showTitle: true,
          color : ['#bdbdbd','#bdbdbd','#bdbdbd','#bdbdbd'],
          feature : {
            mark : {show: false},
            dataZoom : {
              show : true,
              title : {
                dataZoom : 'Data Zoom',
                dataZoomReset : 'Reset Zoom'
              }
            },
            dataView : {show: false, readOnly: true},
            magicType : {
              show: true,
              title : {
                line : 'Line',
                bar : 'Bar',
              },
              type: ['line', 'bar'],
              option: {
                /*line: {
                 itemStyle: {
                 normal: {
                 color:'rgba(3,1,1,1.0)',
                 }
                 },
                 data: [1,2,3,4,5,6,7,8,9,10,11,12]
                 }*/
              }
            },
            restore : {show: false},
            saveAsImage : {show: true,title:'Save as Image'}
          }
        },

        // Enable drag recalculate
        calculable: true,

        // Horizontal axis
        xAxis: [{
          type: 'category',
          boundaryGap: false,
          data: [
            '0h-2h', '2h-4h', '4h-6h', '6h-8h', '8h-10h', '10h-12h', '12h-14h', '14h-16h', '16h-18h', '18h-20h', '20h-22h', '22h-24h'
          ],
          axisLine: {
            show: true,
            onZero: true,
            lineStyle: {
              color: '#ddd',
              type: 'solid',
              width: '1',
              shadowColor: 'rgba(0,0,0,0)',
              shadowBlur: 5,
              shadowOffsetX: 3,
              shadowOffsetY: 3,
            },
          },
          axisTick: {
            show: false,
          },
          splitLine: {
            show: false,
            lineStyle: {
              color: '#fff',
              type: 'solid',
              width: 0,
              shadowColor: 'rgba(0,0,0,0)',
            },
          },
        }],

        // Vertical axis
        yAxis: [{
          type: 'value',
          splitLine: {
            show: false,
            lineStyle: {
              color: 'fff',
              type: 'solid',
              width: 0,
              shadowColor: 'rgba(0,0,0,0)',
            },
          },
          axisLabel: {
            show: false,
          },
          axisTick: {
            show: false,
          },
          axisLine: {
            show: false,
            onZero: true,
            lineStyle: {
              color: '#ff0000',
              type: 'solid',
              width: '0',
              shadowColor: 'rgba(0,0,0,0)',
              shadowBlur: 5,
              shadowOffsetX: 3,
              shadowOffsetY: 3,
            },
          },


        }],

        // Add series
        series: [
          {
            name: 'Page Views',
            type: 'line',
            smooth: true,
            symbol:'none',
            symbolSize:2,
            showAllSymbol: true,
            barWidth:10,
            itemStyle: {
              normal: {
                color:'#64c5b1',
                borderWidth:2, borderColor:'#64c5b1',
                areaStyle: {color:'rgba(100,197,177,0)', type: 'default'}
              }
            },

            data: [1545,1343,1445,2675,2878,1789,1745,2343,2445,1675,1878,2789,1545,1343,1445,2675,2878,1789,1745,2343,2445,1675,1878,2789]
          },
          {
            name: 'Page Views',
            type: 'line',
            smooth: true,
            symbol:'none',
            symbolSize:2,
            showAllSymbol: true,
            barWidth:10,
            itemStyle: {
              normal: {
                color:'#dddddd',
                borderWidth:2, borderColor:'#dddddd)',
                areaStyle: {color:'rgba(221,221,221,0)', type: 'default'}
              }
            },

            data: [5656,6567,7675,3423,4343,5656,6567,7675,3423,4343,5656,6567,7675,3423,4343,5656,6567,7675,3423,4343,5656,6567,7675,3423,4343,5656,6567,7675,3423,4343]
          }
        ]
      };

      // Load data into the ECharts instance
      myChart.setOption(option);

    }




  }



  /******************************
   initialize respective scripts
   *****************************/
  $(document).ready(function() {
    AdminoxAdmin.dashboardEcharts();
  });

  $(window).load(function() {});

});



!function($) {
  "use strict";

  var ChartC3 = function() {};

  ChartC3.prototype.init = function () {

    //Donut Chart
    c3.generate({
      bindto: '#donut-chart',
      data: {
        columns: [
          ['Male', 46],
          ['Female', 24]
        ],
        type : 'donut'
      },
      donut: {
        title: "Candidates",
        width: 30,
        label: {
          show:false
        }
      },
      color: {
        pattern: ["#64c5b1", "#ddd"]
      }
    });

    //Pie Chart
    c3.generate({
      bindto: '#pie-chart',
      data: {
        columns: [
          ['Done', 46],
          ['Due', 24],
          ['Hold', 30]
        ],
        type : 'pie'
      },
      color: {
        pattern: ["#dddddd", "#64c5b1", "#e68900"]
      },
      pie: {
        label: {
          show: false
        }
      }
    });

  },
      $.ChartC3 = new ChartC3, $.ChartC3.Constructor = ChartC3

}(window.jQuery),

//initializing
    function($) {
      "use strict";
      $.ChartC3.init()
    }(window.jQuery);
