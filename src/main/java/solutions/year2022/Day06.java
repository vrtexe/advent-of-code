package solutions.year2022;

import common.Problem;

import java.util.LinkedList;
import java.util.concurrent.atomic.AtomicLong;
import java.util.stream.Stream;

@SuppressWarnings("unused")
public class Day06 extends Problem {

    public static final String file = "2022/06.txt";
    public static final Problem INSTANCE = new Day06(file, SolutionMethod.ALL);

    public Day06(String file) {
        super(file);
    }

    public Day06(String file, SolutionMethod solutionMethod) {
        super(file, solutionMethod);
    }

    public static void main(String[] args) throws Exception {
        INSTANCE.run();
    }

    public String solutionPartOne(Stream<String> data) {
        final var dataStream = data.toList().get(0);
        return String.valueOf(findStartOfPacketMarker(dataStream, 4));
    }

    public String solutionPartTwo(Stream<String> data) {
        final var dataStream = data.toList().get(0);
        return String.valueOf(findStartOfPacketMarker(dataStream, 14));
    }

    public Long findStartOfPacketMarker(String dataStream, Integer uniqueCharacters) {

        final var dataQueue = new LinkedList<String>();
        final var markerLocation = new AtomicLong();

        for (final var character : dataStream.split("")) {
            if (dataQueue.stream().distinct().count() == uniqueCharacters) {
                return markerLocation.get();
            }

            dataQueue.offer(character);
            if (dataQueue.size() > uniqueCharacters) {
                dataQueue.poll();
            }

            markerLocation.incrementAndGet();
        }

        return markerLocation.get();
    }
}
