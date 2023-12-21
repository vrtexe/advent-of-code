package solutions.year2022;

import common.Problem;

import java.util.Comparator;
import java.util.List;
import java.util.Objects;
import java.util.regex.Pattern;
import java.util.stream.LongStream;
import java.util.stream.Stream;

public class Day15 extends Problem {

    public static final String file = "2022/15.txt";
    public static final Problem INSTANCE = new Day15(file);

    /*
     * pos ->
     * pos.distance <=sen.closest.beacon.position
     * cannot have a beacon
     *
     */

    public Day15(String file) {

        super(file);
    }

    public Day15(String file, SolutionMethod solutionMethod) {
        super(file, solutionMethod);
    }

    @Override
    public String solutionPartOne(Stream<String> data) {
        final var y = 2_000_000L;
        // final var y = 10L;

        final var sensors = data.map(this::parseSensor).toList();
        final var sensorsMap = new SensorMap(sensors);

        final var from = sensorsMap.findMinX();
        final var to = sensorsMap.findMaxX();

        final var coveredPositions = LongStream.range(from, to)
                .mapToObj(x -> new Position(x, y))
                .filter(sensorsMap::isCoveringPosition)
                .count();
        // 5838453
        // 5448566
        return String.valueOf(coveredPositions);
    }

    @Override
    public String solutionPartTwo(Stream<String> data) {
        return null;
    }

    public Sensor parseSensor(String line) {
        final var statements = line.split(":");

        final var sensorPosition = extractSensorPosition(statements[0]);
        final var beaconPosition = extractBeaconPosition(statements[1]);

        return new Sensor(sensorPosition, Beacon.of(beaconPosition, sensorPosition));
    }

    public Position extractSensorPosition(String sensorStatement) {
        final var pattern = Pattern.compile("Sensor at x=(-?[0-9]+), y=(-?[0-9]+)");
        final var matcher = pattern.matcher(sensorStatement.trim());

        if (matcher.find()) {
            return new Position(Long.parseLong(matcher.group(1)), Long.parseLong(matcher.group(2)));
        }

        throw new RuntimeException("Invalid sensor string");
    }

    public Position extractBeaconPosition(String beaconStatement) {
        final var pattern = Pattern.compile("closest beacon is at x=(-?[0-9]+), y=(-?[0-9]+)");
        final var matcher = pattern.matcher(beaconStatement.trim());

        if (matcher.find()) {
            return new Position(Long.parseLong(matcher.group(1)), Long.parseLong(matcher.group(2)));
        }

        throw new RuntimeException("Invalid beacon string");
    }

    private record SensorMap(List<Sensor> sensors) {

        public SensorMap(List<Sensor> sensors) {
            this.sensors = sensors.stream()
                    .sorted(Comparator.comparingLong(s -> s.position().x() + s.beacon().distance())).toList();
        }

        public Boolean isCoveringPosition(Position position) {
            return sensors.stream().anyMatch(sensor -> sensor.isCoveringPosition(position));
        }

        public Long findMinX() {
            // return sensors.stream()
            // .min(Comparator.comparing(s -> s.position().x())).stream()
            // .mapToLong(s -> s.position().x() - s.beacon().distance())
            // .findFirst()
            // .orElseThrow();
            final var minSensor = this.sensors().get(0);
            return minSensor.position().x() - minSensor.beacon().distance();
        }

        public Long findMaxX() {
            // return sensors.stream()
            // .max(Comparator.comparing(s -> s.position().x())).stream()
            // .mapToLong(s -> s.position().x() + s.beacon().distance())
            // .findFirst()
            // .orElseThrow();
            final var maxSensor = this.sensors().get(this.sensors().size() - 1);

            return maxSensor.position().x() + maxSensor.beacon().distance();
        }
    }

    private record Sensor(Position position, Beacon beacon) {
        public Boolean isCoveringPosition(Position position) {
            return !beacon.position().equals(position)
                    && this.position().distance(position) <= this.beacon().distance();
        }
    }

    private record Beacon(Position position, Long distance) {
        public static Beacon of(Position position, Position sensor) {
            return new Beacon(position, position.distance(sensor));
        }
    }

    private record Position(Long x, Long y) {
        public Long distance(Position other) {
            final var horizontal = Math.abs(other.x() - this.x());
            final var vertical = Math.abs(other.y() - this.y());

            return Double.valueOf(Math.hypot(vertical, horizontal)).longValue();
        }

        @Override
        public boolean equals(Object other) {
            if (this == other) {
                return true;
            }

            if (other instanceof Position that) {
                return this.x().equals(that.x())
                        && this.y().equals(that.y());
            }

            return true;
        }

        @Override
        public int hashCode() {
            return Objects.hash(this.x(), this.y());
        }
    }

    public static void main(String[] args) throws Exception {
        System.out.println(System.getProperty("platform"));
        INSTANCE.run();
    }
}
