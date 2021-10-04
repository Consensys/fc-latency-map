import React, { useRef, useLayoutEffect } from "react";
import logo from "./logo.svg";
// import "./Map.css";
import * as am4core from "@amcharts/amcharts4/core";
import * as am4maps from "@amcharts/amcharts4/maps";
import am4geodata_worldLow from "@amcharts/amcharts4-geodata/worldLow";
import am4themes_animated from "@amcharts/amcharts4/themes/animated";

am4core.useTheme(am4themes_animated);

interface Props {
  data: any;
}

const Map = (props: Props) => {
  const { data } = props;
  const dataJson = JSON.parse(data);
  const chart = useRef(null);

  /* Add legend */
  function addLegend(chart: any) {
    var legend = new am4maps.Legend();
    legend.parent = chart.chartContainer;
    legend.background.fill = am4core.color("#000");
    legend.background.fillOpacity = 0.05;
    legend.width = 120;
    legend.align = "right";
    legend.padding(10, 15, 10, 15);
    legend.data = [
      {
        name: "Miners",
        fill: "#B27799",
      },
      {
        name: "Probes",
        fill: "#0074d9",
      },
    ];
    legend.itemContainers.template.clickable = false;
    legend.itemContainers.template.focusable = false;

    var legendTitle = legend.createChild(am4core.Label);
    legendTitle.text = "Legend:";
  }

  /* Add miners */
  function addMiners(chart: any, minersList: any) {
    const imageSeries = chart.series.push(new am4maps.MapImageSeries());
    var imageSeriesTemplate = imageSeries.mapImages.template;
    var circle = imageSeriesTemplate.createChild(am4core.Circle);
    circle.radius = 4;
    circle.fill = am4core.color("#B27799");
    circle.stroke = am4core.color("#FFFFFF");
    circle.strokeWidth = 2;
    circle.nonScaling = true;
    circle.tooltipText = "{title}";

    imageSeriesTemplate.propertyFields.latitude = "latitude";
    imageSeriesTemplate.propertyFields.longitude = "longitude";

    // const test = minersList.filter((info: any) => {
    //   if (info.latitude && info.longitude) {
    //     return {
    //       latitude: info.latitude,
    //       longitude: info.longitude,
    //       title: `${info.address}\n${info.ip}`,
    //     };
    //   }
    //   return null;
    // });

    // const test = minersList.filter((info: any) => {
    //   if (info.latitude && info.longitude) {
    //     return {
    // latitude: info.latitude,
    // longitude: info.longitude,
    // title: `${info.address}\n${info.ip}`,
    //     };
    //   }
    //   return null;
    // });

    const test = minersList.reduce((a, info) => {
      info.latitude &&
        info.longitude &&
        a.push({
          latitude: info.latitude,
          longitude: info.longitude,
          title: `${info.address}\n${info.ip}`,
        });
      return a;
    }, []);

    console.log("xxx==>>", test);

    imageSeries.data = test;
  }

  /* Add miners */
  function addMiners3(chart: any, minersList: any) {
    var imageSeries = chart.series.push(new am4maps.MapImageSeries());
    imageSeries.mapImages.template.propertyFields.longitude = "longitude";
    imageSeries.mapImages.template.propertyFields.latitude = "latitude";
    imageSeries.mapImages.template.tooltipText = "{title}";
    imageSeries.mapImages.template.propertyFields.url = "url";

    var circle = imageSeries.mapImages.template.createChild(am4core.Circle);
    circle.radius = 3;
    circle.propertyFields.fill = "color";
    circle.nonScaling = true;

    var circle2 = imageSeries.mapImages.template.createChild(am4core.Circle);
    circle2.radius = 3;
    circle2.propertyFields.fill = "color";

    circle2.events.on("inited", function (event) {
      animateBullet(chart, event.target);
    });

    // imageSeriesTemplate.propertyFields.latitude = "latitude";
    // imageSeriesTemplate.propertyFields.longitude = "longitude";

    // imageSeries.data = minersList.filter((info: any) => {
    //   if (info.latitude && info.longitude) {
    //     return {
    //       latitude: info.latitude,
    //       longitude: info.longitude,
    //       title: `${info.address}\n${info.ip}`,
    //     };
    //   }
    // });
    imageSeries.data = [
      {
        latitude: 59.9452,
        longitude: 10.7559,
        title: `f023467`,
        color: "#ffff00",
      },
      {
        latitude: 52.48395,
        longitude: -1.8898,
        title: `f0694396`,
        color: "#ffff00",
      },
      {
        latitude: 47.36329,
        longitude: 8.55014,
        title: `f022163`,
        color: "#00ff00",
      },
      {
        latitude: 37.55983,
        longitude: -122.27148,
        title: `f01231`,
        color: "#ff0000",
      },
      {
        latitude: 37.41043,
        longitude: 127.13716,
        title: `f01044351`,
        color: "#ff0000",
      },
    ];
  }

  function animateBullet(chart, circle) {
    var animation = circle.animate(
      [
        {
          property: "scale",
          from: 1 / chart.zoomLevel,
          to: 5 / chart.zoomLevel,
        },
        { property: "opacity", from: 1, to: 0 },
      ],
      1000,
      am4core.ease.circleOut
    );
    animation.events.on("animationended", function (event) {
      animateBullet(chart, event.target.object);
    });
  }

  /* Add probes */
  function addProbes(chart: any, color: string) {
    const imageSeries = chart.series.push(new am4maps.MapImageSeries());
    var imageSeriesTemplate = imageSeries.mapImages.template;
    var circle2 = imageSeriesTemplate.createChild(am4core.Circle);
    circle2.radius = 4;
    circle2.fill = am4core.color(color);
    // circle2.fill = am4core.color("#0074d9");
    circle2.stroke = am4core.color("#FFFFFF");
    circle2.strokeWidth = 2;
    circle2.nonScaling = true;
    circle2.tooltipText = "{title}";

    imageSeriesTemplate.propertyFields.latitude = "latitude";
    imageSeriesTemplate.propertyFields.longitude = "longitude";
    imageSeries.data = [
      {
        latitude: 49.9452,
        longitude: 10.7559,
        title: "Probe",
      },
      {
        latitude: 9.9452,
        longitude: 22.7559,
        title: "Probe",
      },
    ];

    imageSeries.events.on("hit", function (ev) {
      // get object info
      console.log(ev.target);
      addMiners3(chart, dataJson["CN"]["PEK"]);
      // toto.mapLines.clear();
      // ev.target.fill = am4core.color("#ff00ff");
      // let series = ev.target.chart.series.push(new am4maps.ColumnSeries());
    });
  }

  useLayoutEffect(() => {
    // Map config
    let chart = am4core.create("chartdiv", am4maps.MapChart);
    chart.projection = new am4maps.projections.Miller();

    chart.geodata = am4geodata_worldLow;

    let polygonSeries = new am4maps.MapPolygonSeries();
    polygonSeries.useGeodata = true;
    chart.series.push(polygonSeries);

    // Remove Antarctica
    polygonSeries.exclude = ["AQ"];

    // Configure series
    var polygonTemplate = polygonSeries.mapPolygons.template;
    polygonTemplate.tooltipText = "{name}";
    polygonTemplate.polygon.fillOpacity = 0.6;

    // Create hover state and set alternative fill color
    var hs = polygonTemplate.states.create("hover");
    hs.properties.fill = chart.colors.getIndex(0);

    // Lines
    // let lineSeries = chart.series.push(new am4maps.MapLineSeries());
    // lineSeries.data = [
    //   {
    //     multiGeoLine: [
    //       [
    //         { latitude: 48.856614, longitude: 2.352222 },
    //         { latitude: 40.712775, longitude: -74.005973 },
    //         { latitude: 49.282729, longitude: -123.120738 },
    //       ],
    //     ],
    //   },
    // ];

    // Legend
    addLegend(chart);

    // Miners
    addMiners(chart, dataJson["CN"]["PEK"]);

    // Probes
    addProbes(chart, "#0074d9");

    // Events
    // var polygonTemplate = polygonSeries.mapPolygons.template;
    // polygonTemplate.events.on("hit", function (ev) {
    //   // zoom to an object
    //   ev.target.series.chart.zoomToMapObject(ev.target);

    //   // get object info
    //   console.log(ev.target.dataItem.dataContext.name);
    // });

    return () => {
      chart.dispose();
    };
  }, []);

  function minerClickHandler(ev) {
    console.log("clicked on ", ev.target);

    ev.target.series.chart.zoomToMapObject(ev.target);
    lineSeries.mapLines.clear();
  }

  return <div id="chartdiv" style={{ width: "100%", height: "1000px" }}></div>;
  // return <div id="chartdiv" style={{ width: "1000px", height: "500px" }}></div>;
};

export default Map;
