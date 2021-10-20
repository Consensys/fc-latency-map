import React, { useRef, useState, useLayoutEffect } from "react";
import { Row, Col } from "antd";
import getConfig from "next/config";
import * as am4core from "@amcharts/amcharts4/core";
import * as am4maps from "@amcharts/amcharts4/maps";
import am4geodata_worldLow from "@amcharts/amcharts4-geodata/worldLow";
import am4themes_animated from "@amcharts/amcharts4/themes/animated";
import { Button } from "antd";
import { LeftOutlined, RightOutlined } from "@ant-design/icons";

import styles from "@src/styles/Global.module.css";
import Location from "./Location";
import Miner from "./Miner";

const { publicRuntimeConfig } = getConfig();

am4core.useTheme(am4themes_animated);

interface Miner {
  address: string;
  ip: string;
  port: string;
  latitude: number;
  longitude: number;
}

interface Location {
  name: string;
  country: string;
  iata_code: string;
  latitude: number;
  longitude: number;
  type: string;
}

interface MinerLatency {
  address: string;
  ip: string;
  port: string;
  latitude: number;
  longitude: number;
  latency: {
    avg: number;
  };
}

interface MinerLatencyData {
  address: string;
  measures: LatencyMeasureData[];
}

interface LatencyMeasureData {
  ip: string;
  latency: [];
}

interface Props {
  data: any;
  date: string;
  dates: string[];
  height: any;
  width: any;
}

const Map = (props: Props) => {
  const { data, date, dates, width } = props;

  const dataJson = JSON.parse(data);
  const locations = dataJson.locations ? dataJson.locations : [];
  const miners = dataJson.miners
    ? dataJson.miners.filter((miner: any) => miner.latitude && miner.longitude)
    : [];

  const chart = useRef(null);

  const [location, setLocation] = useState(null);
  const [miner, setMiner] = useState(null);

  /* Add legend */
  function addLegend(chart: any) {
    var legend = new am4maps.Legend();
    legend.parent = chart.chartContainer;
    legend.background.fill = am4core.color("#000");
    legend.background.fillOpacity = 0.05;
    legend.width = 230;
    legend.align = "right";
    legend.fontSize = 11;
    legend.padding(5, 10, 5, 10);
    legend.data = [
      {
        name: "Locations",
        fill: "#C0C0C0",
      },
      {
        name: "Miners",
        fill: "#4169E1",
      },
      {
        name: `Low Latency Miners <= ${publicRuntimeConfig.latency.low}`,
        fill: "#00FF00",
      },
      {
        name: `Medium Latency Miners <= ${publicRuntimeConfig.latency.medium}`,
        fill: "#FFFF00",
      },
      {
        name: `High Latency Miners > ${publicRuntimeConfig.latency.medium}`,
        fill: "#FF0000",
      },
      {
        name: "Not responding Miners",
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

    imageSeries.data = minersList.map((miner: Miner) => {
      return {
        title: miner.address,
        ip: miner.ip,
        port: miner.port,
        latitude: miner.latitude,
        longitude: miner.longitude,
        color: "#4169E1",
      };
    });

    imageSeries.mapImages.template.events.on(`hit`, (e: any) => {
      const miner = e.target.dataItem.dataContext;
      setMiner(miner);
    });
  }

  /* Add miners with latency*/
  function addMinersLatency(chart: any, minersList: any) {
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

    circle2.events.on("inited", function (event: any) {
      animateBullet(chart, event.target);
    });

    imageSeries.data = minersList.map((miner: MinerLatency) => {
      let color = "#ff0000";

      if (miner.latency.avg == -1) {
        color = "#000000";
      } else if (miner.latency.avg < publicRuntimeConfig.latency.low) {
        color = "#00ff00";
      } else if (miner.latency.avg < publicRuntimeConfig.latency.medium) {
        color = "#ffff00";
      } else {
        color = "#ff0000";
      }

      return {
        title: miner.address,
        ip: miner.ip,
        port: miner.port,
        latitude: miner.latitude,
        longitude: miner.longitude,
        color,
        latency: miner.latency,
      };
    });

    imageSeries.mapImages.template.events.on(`hit`, (event: any) => {
      const miner = event.target.dataItem.dataContext;
      setMiner(miner);
    });
  }

  function animateBullet(chart: am4maps.MapChart, circle: any) {
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
    animation.events.on("animationended", (event: any) => {
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

    series.data = locations.map((location: Location) => {
      return {
        title: `${location.iata_code} - ${location.name}`,
        name: location.name,
        iataCode: location.iata_code,
        country: location.country,
        latitude: location.latitude,
        longitude: location.longitude,
      };
    });

    // Location clicked
    series.mapImages.template.events.on(`hit`, (ev: any) => {
      const location = ev.target.dataItem.dataContext;
      setLocation(location);

      const latenciesList =
        dataJson.measurements[location.country] &&
        dataJson.measurements[location.country][location.iataCode]
          ? dataJson.measurements[location.country][location.iataCode]
          : [];

      let minersLatency: MinerLatency[] = [];
      let minersNoLatency: MinerLatency[] = [];

      miners.forEach((miner: Miner, index: number) => {
        const existsLatency = latenciesList.find(
          (latency: MinerLatencyData) => latency.address == miner.address
        );
        if (existsLatency) {
          const minerLatency = {
            ...miners[index],
            latency: existsLatency.measures[0],
          };
          minersLatency.push(minerLatency);
        } else {
          minersNoLatency.push(miners[index]);
        }
      });

      while (chart.series.length >= 3) {
        chart.series.removeIndex(chart.series.length - 1);
      }

      addMinersLatency(chart, minersLatency);
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

  const isPrevious = () => {
    let index = dates.findIndex((dateElt) => dateElt == date);

    return (
      <Button
        shape="circle"
        type="link"
        icon={<LeftOutlined />}
        href={`/date/${dates[index - 1] && dates[index - 1]}`}
        className={styles.dateButton}
        disabled={index == 0}
      />
    );
  };

  const isNext = () => {
    let index = dates.findIndex((dateElt) => dateElt == date);

    return (
      <Button
        shape="circle"
        type="link"
        icon={<RightOutlined />}
        href={`/date/${dates[index + 1] && dates[index + 1]}`}
        className={styles.dateButton}
        disabled={index + 1 == dates.length}
      />
    );
  };

  return (
    <div>
      <div
        id="chartdiv"
        style={{
          width: "100%",
          height: width / 1.7,
          minHeight: "500px",
        }}
      ></div>
      <div className={styles.dates}>
        {isPrevious()}
        <div className={styles.date}>{date}</div>
        {isNext()}
      </div>

      <Row gutter={[16, 16]} className={styles.informations}>
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
