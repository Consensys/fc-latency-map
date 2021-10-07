import { Card } from "antd";
import * as am4core from "@amcharts/amcharts4/core";
import am4themes_animated from "@amcharts/amcharts4/themes/animated";

import styles from "@src/styles/Global.module.css";

am4core.useTheme(am4themes_animated);

interface Props {
  location: any;
}

const Location = (props: Props) => {
  const { location } = props;

  const displayLocation = (location: any) =>
    location ? (
      <Card title="Location" className={styles.information}>
        <p>Name: {location.name}</p>
        <p>Iata Code: {location.iataCode}</p>
        <p>Latitude: {location.latitude}</p>
        <p>Longitude: {location.longitude}</p>
      </Card>
    ) : (
      <Card title="Location" className={styles.information}>
        Please choose a location...
      </Card>
    );

  return <div>{displayLocation(location)}</div>;
};

export default Location;
