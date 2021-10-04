import React, { useRef, useState, useLayoutEffect } from "react";
import { Row, Col } from "antd";

// import logo from "./logo.svg";
// import "./Map.css";
import * as am4core from "@amcharts/amcharts4/core";
import * as am4maps from "@amcharts/amcharts4/maps";
import am4geodata_worldLow from "@amcharts/amcharts4-geodata/worldLow";
import am4themes_animated from "@amcharts/amcharts4/themes/animated";

import Location from "./Location";
import Miner from "./Miner";

am4core.useTheme(am4themes_animated);

interface Props {
  data: any;
  locationsJson: any;
}

const Map = (props: Props) => {
  const { data, locationsJson } = props;
  const dataJson = JSON.parse(data);
  const locations = dataJson.location;
  const miners = dataJson.miners;

  const chart = useRef(null);

  const [location, setLocation] = useState(null);
  const [miner, setMiner] = useState(null);

  /* Add legend */
  function addLegend(chart: any) {
    var legend = new am4maps.Legend();
    legend.parent = chart.chartContainer;
    legend.background.fill = am4core.color("#000");
    legend.background.fillOpacity = 0.05;
    legend.width = 190;
    legend.align = "right";
    legend.fontSize = 11;
    legend.padding(5, 10, 5, 10);
    legend.data = [
      {
        name: "Probes",
        fill: "#C0C0C0",
      },
      {
        name: "Miners",
        fill: "#6F00FF",
      },
      {
        name: "Low Latency Miners",
        fill: "#00FF00",
      },
      {
        name: "Medium Latency Miners",
        fill: "#FFFF00",
      },
      {
        name: "High Latency Miners",
        fill: "#FF0000",
      },
      {
        name: "Timeout Miners",
        fill: "#000000",
      },
    ];
    legend.itemContainers.template.clickable = false;
    legend.itemContainers.template.focusable = false;

    var legendTitle = legend.createChild(am4core.Label);
    legendTitle.text = "Legend:";
  }

  /* Add miners */
  function addMiners(chart: any, minersList: any) {
    var imageSeries = chart.series.push(new am4maps.MapImageSeries());
    imageSeries.mapImages.template.propertyFields.longitude = "longitude";
    imageSeries.mapImages.template.propertyFields.latitude = "latitude";
    imageSeries.mapImages.template.tooltipText = "{title}";
    imageSeries.mapImages.template.propertyFields.url = "url";

    var circle = imageSeries.mapImages.template.createChild(am4core.Circle);
    circle.radius = 3;
    circle.propertyFields.fill = "color";
    circle.nonScaling = true;

    imageSeries.data = minersList.map((miner) => {
      return {
        latitude: miner.latitude,
        longitude: miner.longitude,
        title: miner.address,
        color: "#6F00FF",
      };
    });

    imageSeries.mapImages.template.events.on(`hit`, (ev) => {
      const miner = ev.target.dataItem.dataContext;
      setMiner(miner);
    });
  }

  /* Add miners */
  function addMinersLantecy(chart: any, minersList: any) {
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

    // imageSeries.data = minersList;
    imageSeries.data = minersList.map((miner) => {
      let color = "#ff0000";

      if (miner.latency.avg == -1) {
        color = "#000000";
      } else if (miner.latency.avg < 80) {
        color = "#00ff00";
      } else if (miner.latency.avg < 100) {
        color = "#00ffff";
      } else {
        color = "#ff0000";
      }

      return {
        latitude: miner.latitude,
        longitude: miner.longitude,
        title: miner.address,
        color,
        latency: miner.latency,
      };
    });

    imageSeries.mapImages.template.events.on(`hit`, (ev) => {
      const miner = ev.target.dataItem.dataContext;
      setMiner(miner);
    });
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
    const series = chart.series.push(new am4maps.MapImageSeries());

    var template = series.mapImages.template;
    template.verticalCenter = "middle";
    template.horizontalCenter = "middle";
    template.propertyFields.latitude = "lat";
    template.propertyFields.longitude = "long";
    template.propertyFields.latitude = "latitude";
    template.propertyFields.longitude = "longitude";

    var circle = template.createChild(am4core.Circle);
    circle.radius = 4;
    circle.fill = am4core.color(color);
    circle.stroke = am4core.color("#FFFFFF");
    circle.strokeWidth = 2;
    circle.nonScaling = true;
    circle.tooltipText = "{title}";

    series.data = locations.map((location) => {
      return {
        latitude: location.latitude,
        longitude: location.longitude,
        title: location.iata_code,
        iataCode: location.iata_code,
        country: location.country,
      };
    });

    series.mapImages.template.events.on(`hit`, (ev) => {
      const location = ev.target.dataItem.dataContext;
      setLocation(location);

      const latenciesList =
        dataJson.measurements[location.country][location.iataCode];

      let minersLatency = [];
      let minersNoLatency = [];

      latenciesList.forEach((latency) => {
        const index = miners.findIndex(
          (miner) => miner.address == latency.address
        );

        if (
          latency.measures &&
          !minersLatency.find((miner) => miner.address == latency.address)
        ) {
          const minerLatency = {
            ...miners[index],
            latency: latency.measures[0].latency[0],
          };
          minersLatency.push(minerLatency);
        } else if (
          !latency.measures &&
          !minersNoLatency.find((miner) => miner.address == latency.address)
        ) {
          minersNoLatency.push(miners[index]);
        }
      });

      while (chart.series.length >= 3) {
        chart.series.removeIndex(chart.series.length - 1);
      }

      addMinersLantecy(chart, minersLatency);
      addMiners(chart, minersNoLatency);
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

    // polygonTemplate.fill = am4core.color("#454a58");

    // Create hover state and set alternative fill color
    var hs = polygonTemplate.states.create("hover");
    hs.properties.fill = chart.colors.getIndex(0);

    // Legend
    addLegend(chart);

    // Probes
    addProbes(chart, "#C0C0C0");
    addMiners(chart, miners);

    return () => {
      chart.dispose();
    };
  }, []);

  return (
    <div>
      <div
        id="chartdiv"
        style={{
          width: "100%",
          height: "1000px",
          // backgroundColor: "#454a58"
        }}
      ></div>
      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Location location={location} />
        </Col>
        <Col span={12}>
          <Miner miner={miner} />
        </Col>
      </Row>
    </div>
  );
};

export default Map;
