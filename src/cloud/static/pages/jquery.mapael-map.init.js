/**
 * Theme: Adminox Admin Template
 * Author: Coderthemes
 * Module/App: Mapael Maps
 */


$(function(){

	//USA Map

	$mapusa = $(".map-usa");

 	$mapusa.mapael({
		map : {
			name : "usa_states",
             defaultArea: {
                attrs: {
                    fill: "#36404e",
                    stroke: "#aaa"
                },
                 attrsHover: {
                    fill: "#4489e4"
                }
            },
			zoom: {
				enabled: true,
				maxLevel : 10
			}
		},
        legend: {
            plot: {
                title: "American cities",
                slices: [{
                    size: 24,
                    attrs: {
                        fill: "#188ae2"
                    },
                    label: "Product One",
                    sliceValue: "Value 1"
                }, {
                    size: 24,
                    attrs: {
                        fill: "#3ac9d6"
                    },
                    label: "Product Two",
                    sliceValue: "Value 2"
                }, {
                    size: 24,
                    attrs: {
                        fill: "#f5707a"
                    },
                    label: "Product Three",
                    sliceValue: "Value 3"
                }]
            }
        },
		plots: {
			'ny': {
                latitude: 40.717079,
                longitude: -74.00116,
                tooltip: {content: "New York"},
                value: "Value 3"
            },
            'an': {
                latitude: 61.2108398,
                longitude: -149.9019557,
                tooltip: {content: "Anchorage"},
                value: "Value 3"
            },
            'sf': {
                latitude: 37.792032,
                longitude: -122.394613,
                tooltip: {content: "San Francisco"},
                value: "Value 1"
            },
            'pa': {
                latitude: 19.493204,
                longitude: -154.8199569,
                tooltip: {content: "Pahoa"},
                value: "Value 2"
            },
            'la': {
                latitude: 34.025052,
                longitude: -118.192006,
                tooltip: {content: "Los Angeles"},
                value: "Value 3"
            },
            'dallas': {
                latitude: 32.784881,
                longitude: -96.808244,
                tooltip: {content: "Dallas"},
                value: "Value 2"
            },
            'miami': {
                latitude: 25.789125,
                longitude: -80.205674,
                tooltip: {content: "Miami"},
                value: "Value 3"
            },
            'washington': {
                latitude: 38.905761,
                longitude: -77.020746,
                tooltip: {content: "Washington"},
                value: "Value 2"
            },
            'seattle': {
                latitude: 47.599571,
                longitude: -122.319426,
                tooltip: {content: "Seattle"},
                value: "Value 1"
            }
		}
	});

	// Zoom on mousewheel with mousewheel jQuery plugin
	$mapusa.on("mousewheel", function(e) {
		if (e.deltaY > 0) {
			$mapusa.trigger("zoom", $mapusa.data("zoomLevel") + 1);
			console.log("zoom");
		} else {
			$mapusa.trigger("zoom", $mapusa.data("zoomLevel") - 1);
		}

		return false;
	});


    $(".mapcontainer").mapael({
        map: {
            name: "world_countries",
            defaultArea: {
                attrs: {
                    fill: "#36404e"
                    , stroke: "#aaa"
                },
                attrsHover: {
                    fill: "#ff9800"
                }
            }
            // Default attributes can be set for all links
            , defaultLink: {
                factor: 0.4
                , attrsHover: {
                    stroke: "#f06292"
                }
            }
            , defaultPlot: {
                text: {
                    attrs: {
                        fill: "#ddd"
                    },
                    attrsHover: {
                        fill: "#ddd"
                    }
                }
            }
        },
        plots: {
            'paris': {
                latitude: 48.86,
                longitude: 2.3444,
                tooltip: {content: "Paris<br />Population: 500000000"}
            },
            'newyork': {
                latitude: 40.667,
                longitude: -73.833,
                tooltip: {content: "New york<br />Population: 200001"}
            },
            'sanfrancisco': {
                latitude: 37.792032,
                longitude: -122.394613,
                tooltip: {content: "San Francisco"}
            },
            'brasilia': {
                latitude: -15.781682,
                longitude: -47.924195,
                tooltip: {content: "Brasilia<br />Population: 200000001"}
            },
            'roma': {
                latitude: 41.827637,
                longitude: 12.462732,
                tooltip: {content: "Roma"}
            },
            'miami': {
                latitude: 25.789125,
                longitude: -80.205674,
                tooltip: {content: "Miami"}
            },

            // Size=0 in order to make plots invisible
            'tokyo': {
                latitude: 35.687418,
                longitude: 139.692306,
                size: 0,
                text: {content: 'Tokyo'}
            },
            'sydney': {
                latitude: -33.917,
                longitude: 151.167,
                size: 0,
                text: {content: 'Sydney'}
            },
            'plot1': {
                latitude: 22.906561,
                longitude: 86.840170,
                size: 0,
                text: {content: 'Plot1', position: 'left', margin: 5}
            },
            'plot2': {
                latitude: -0.390553,
                longitude: 115.586762,
                size: 0,
                text: {content: 'Plot2'}
            },
            'plot3': {
                latitude: 44.065626,
                longitude: 94.576079,
                size: 0,
                text: {content: 'Plot3'}
            }
        },
        // Links allow you to connect plots between them
        links: {
            'link1': {
                factor: -0.3
                // The source and the destination of the link can be set with a latitude and a longitude or a x and a y ...
                , between: [{latitude: 24.708785, longitude: -5.402427}, {x: 560, y: 280}]
                , attrs: {
                    "stroke-width": 2
                }
                , tooltip: {content: "Link"}
            }
            , 'parisnewyork': {
                // ... Or with IDs of plotted points
                factor: -0.3
                , between: ['paris', 'newyork']
                , attrs: {
                    "stroke-width": 2
                }
                , tooltip: {content: "Paris - New-York"}
            }
            , 'parissanfrancisco': {
                // The curve can be inverted by setting a negative factor
                factor: -0.5
                , between: ['paris', 'sanfrancisco']
                , attrs: {
                    "stroke-width": 4
                }
                , tooltip: {content: "Paris - San - Francisco"}
            }
            , 'parisbrasilia': {
                factor: -0.8
                , between: ['paris', 'brasilia']
                , attrs: {
                    "stroke-width": 1
                }
                , tooltip: {content: "Paris - Brasilia"}
            }
            , 'romamiami': {
                factor: 0.2
                , between: ['roma', 'miami']
                , attrs: {
                    "stroke-width": 4
                }
                , tooltip: {content: "Roma - Miami"}
            }
            , 'sydneyplot1': {
                factor: -0.2
                , between: ['sydney', 'plot1']
                , attrs: {
                    stroke: "#4489e4",
                    "stroke-width": 3,
                    "stroke-linecap": "round",
                    opacity: 0.6
                }
                , tooltip: {content: "Sydney - Plot1"}
            }
            , 'sydneyplot2': {
                factor: -0.1
                , between: ['sydney', 'plot2']
                , attrs: {
                    stroke: "#4489e4",
                    "stroke-width": 8,
                    "stroke-linecap": "round",
                    opacity: 0.6
                }
                , tooltip: {content: "Sydney - Plot2"}
            }
            , 'sydneyplot3': {
                factor: 0.2
                , between: ['sydney', 'plot3']
                , attrs: {
                    stroke: "#4489e4",
                    "stroke-width": 4,
                    "stroke-linecap": "round",
                    opacity: 0.6
                }
                , tooltip: {content: "Sydney - Plot3"}
            }
            , 'sydneytokyo': {
                factor: 0.2
                , between: ['sydney', 'tokyo']
                , attrs: {
                    stroke: "#4489e4",
                    "stroke-width": 6,
                    "stroke-linecap": "round",
                    opacity: 0.6
                }
                , tooltip: {content: "Sydney - Plot2"}
            }
        }
    });

});