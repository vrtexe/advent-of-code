package solutions.year2022;

import common.Problem;

import java.util.List;
import java.util.Map;
import java.util.stream.Stream;


@SuppressWarnings("unused")
public class Day02 extends Problem {

    public static final List<String> MyMoves = List.of("A", "B", "C");
    private static final String file = "2022/02.txt";
    private static final Problem INSTANCE = new Day02(file);

    public Day02(String file) {
        super(file);
    }

    public Day02(String file, SolutionMethod solutionMethod) {
        super(file, solutionMethod);
    }

    public static void main(String[] args) throws Exception {
        INSTANCE.run();
    }

    @Override
    public String solutionPartOne(Stream<String> data) {
        var totalPoints = data.map(this::getMoveEvaluator).mapToLong(MoveEvaluator::eval).sum();
        return String.valueOf(totalPoints);
    }

    @Override
    public String solutionPartTwo(Stream<String> data) {
        var totalPoints = data.map(this::getMoveEvaluatorAlt).mapToLong(MoveEvaluator::eval).sum();

        return String.valueOf(totalPoints);
    }

    private MoveEvaluator getMoveEvaluator(String round) {
        var moves = round.split(" ");
        return new MoveEvaluator(moves[1], moves[0]);
    }

    private MoveEvaluator getMoveEvaluatorAlt(String round) {
        var moves = round.split(" ");
        return new MoveEvaluator(moves[1], moves[0], true);
    }

    private enum Moves {
        ROCK("ROCK"), PAPER("PAPER"), SCISSORS("SCISSORS");

        private final String value;

        Moves(String value) {
            this.value = value;
        }

        public String getValue() {
            return value;
        }
    }

    private static class MoveEvaluator {
        public static final Map<String, Moves> OpponentMoves = Map.of("A", Moves.ROCK, "B", Moves.PAPER, "C",
                Moves.SCISSORS);
        public static final Map<String, Moves> MyMoves = Map.of("X", Moves.ROCK, "Y", Moves.PAPER, "Z", Moves.SCISSORS);

        public static final Map<Moves, Moves> LosingMove = Map.of(Moves.PAPER, Moves.ROCK, Moves.SCISSORS, Moves.PAPER,
                Moves.ROCK, Moves.SCISSORS);

        public static final Map<Moves, Moves> WinningMove = Map.of(Moves.ROCK, Moves.PAPER, Moves.SCISSORS, Moves.ROCK,
                Moves.PAPER, Moves.SCISSORS);

        private Moves myMove;
        private Moves opponentMove;

        public MoveEvaluator(String myMove, String opponentMove) {
            createRegularMode(myMove, opponentMove);
        }

        public MoveEvaluator(String myMove, String opponentMove, Boolean alt) {
            if (!alt) {
                createRegularMode(myMove, opponentMove);
            } else {
                createDictationMode(myMove, opponentMove);
            }
        }

        private void createRegularMode(String myMove, String opponentMove) {
            this.myMove = MyMoves.get(myMove);
            this.opponentMove = OpponentMoves.get(opponentMove);
        }

        private void createDictationMode(String myMove, String opponentMove) {
            this.opponentMove = OpponentMoves.get(opponentMove);
            this.myMove = getAltMoveMapping(myMove, this.opponentMove);
        }

        private Moves getAltMoveMapping(String myMove, Moves opponentMove) {
            return switch (myMove) {
                case "X" -> LosingMove.get(opponentMove);
                case "Y" -> opponentMove;
                case "Z" -> WinningMove.get(opponentMove);
                default -> opponentMove;
            };
        }

        public Long eval() {
            return getShapeScore(myMove) + getGameScore(myMove, opponentMove);
        }

        private Long getGameScore(Moves myMove, Moves opponentMove) {
            if (myMove == opponentMove) {
                return 3L;
            }

            if (isLoss(myMove, opponentMove)) {
                return 0L;
            }

            if (isWin(myMove, opponentMove)) {
                return 6L;
            }

            return 0L;

        }

        private Boolean isLoss(Moves myMove, Moves opponentMove) {
            return myMove == Moves.SCISSORS && opponentMove == Moves.ROCK
                    || myMove == Moves.ROCK && opponentMove == Moves.PAPER
                    || myMove == Moves.PAPER && opponentMove == Moves.SCISSORS;
        }

        private Boolean isWin(Moves myMove, Moves opponentMove) {
            return opponentMove == Moves.SCISSORS && myMove == Moves.ROCK
                    || opponentMove == Moves.ROCK && myMove == Moves.PAPER
                    || opponentMove == Moves.PAPER && myMove == Moves.SCISSORS;
        }

        private Long getShapeScore(Moves myMove) {
            return switch (myMove) {
                case ROCK -> 1L;
                case PAPER -> 2L;
                case SCISSORS -> 3L;
            };
        }

        @Override
        public String toString() {
            return "%s - %s".formatted(myMove.getValue(), opponentMove.getValue());
        }
    }
}
