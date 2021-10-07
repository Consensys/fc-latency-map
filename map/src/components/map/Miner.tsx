import { Card } from "antd";
import styles from "@src/styles/Global.module.css";

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
      <Card title="Miner" className={styles.information}>
        <p>Name: {miner.title}</p>
        <p>Latitude: {miner.latitude}</p>
        <p>Longitude: {miner.longitude}</p>

        <p>
          Latency:
          {miner.latency ? (
            <ul>
              <li>Average: {miner.latency.avg}</li>
              <li>Max: {miner.latency.max}</li>
              <li>Min: {miner.latency.min}</li>
            </ul>
          ) : (
            <> N/A</>
          )}
        </p>
      </Card>
    );

  return <div>{miner && displayMiner(miner)}</div>;
};

export default Miner;
