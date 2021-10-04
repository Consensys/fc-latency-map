import React, { useRef, useState, useLayoutEffect } from "react";
import { Card } from "antd";
import { UpCircleOutlined } from "@ant-design/icons";

import * as am4core from "@amcharts/amcharts4/core";
import am4themes_animated from "@amcharts/amcharts4/themes/animated";

am4core.useTheme(am4themes_animated);

interface Props {
  miner: any;
}

const Miner = (props: Props) => {
  const { miner } = props;

  const displayMiner = (miner: any) =>
    miner && (
      <Card title="Miner">
        <p>Name: {miner.title}</p>
        <p>Latitude: {miner.latitude}</p>
        <p>Longitude: {miner.longitude}</p>
        {miner.latency && (
          <p>
            Latency:
            <ul>
              <li>Average: {miner.latency.avg}</li>
              <li>Max: {miner.latency.max}</li>
              <li>Min: {miner.latency.min}</li>
            </ul>
          </p>
        )}
      </Card>
    );

  return <div>{miner && displayMiner(miner)}</div>;
};

export default Miner;
