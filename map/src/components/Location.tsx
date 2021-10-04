import React, { useRef, useState, useLayoutEffect } from "react";
import { Card } from "antd";

import * as am4core from "@amcharts/amcharts4/core";
import am4themes_animated from "@amcharts/amcharts4/themes/animated";

am4core.useTheme(am4themes_animated);

interface Props {
  location: any;
}

const Location = (props: Props) => {
  const { location } = props;

  const displayLocation = (location: any) =>
    location && (
      <Card title="Location">
        <p>Name: {location.title}</p>
        <p>Latitude: {location.latitude}</p>
        <p>Longitude: {location.longitude}</p>
      </Card>
    );

  return <div>{location && displayLocation(location)}</div>;
};

export default Location;
