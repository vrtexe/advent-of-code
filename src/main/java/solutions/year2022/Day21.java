package solutions.year2022;

import common.Problem;

import java.util.*;
import java.util.stream.Collectors;
import java.util.stream.Stream;

public class Day21 extends Problem {

    public static final String file = "2022/21.txt";
    public static final Problem INSTANCE = new Day21();

    public static final List<String> SUPPORTED_OPERATIONS = List.of("+", "-", "");

    @Override
    public String solutionPartOne(Stream<String> data) {
        var monkeyMap = generateData(data);

        return getSolutionPart1(monkeyMap);
    }

    @Override
    public String solutionPartTwo(Stream<String> data) {
        var monkeyMap = generateData(data);

        return getSolutionPart2(monkeyMap);
    }

    private Map<String, String> generateData(Stream<String> data) {
        return generateData(data, false);
    }

    private Map<String, String> generateData(Stream<String> data, Boolean stub) {
        if (stub) {
            return generateStubData();
        }

        return mapData(data);
    }

    private Map<String, String> mapData(Stream<String> data) {
        return data.map(monkey -> monkey.split(":"))
                .collect(Collectors.toMap(monkey -> monkey[0].trim(), monkey -> monkey[1].trim()));
    }

    private Map<String, String> generateStubData() {
        var stubData = List.of("root: pppw + sjmn", "dbpl: 5", "cczh: sllz + lgvd", "zczc: 2", "ptdq: humn - dvpt",
                "dvpt: 3", "lfqf: 4", "humn: 5", "ljgn: 2", "sjmn: drzm * dbpl", "sllz: 4",
                "pppw: cczh / lfqf", "lgvd: ljgn * ptdq", "drzm: hmdt - zczc", "hmdt: 32").stream();

        return mapData(stubData);
    }

    private String getSolutionPart1(Map<String, String> monkeyMap) {
        return eval("root", monkeyMap.get("root"), monkeyMap).toString();
    }

    private String getSolutionPart2(Map<String, String> monkeyMap) {
        return eval(monkeyMap).toString();
    }

    private Long eval(String monkeyName, String rootEXpression, Map<String, String> monkeyMap) {

        final var expressionStack = new Stack<Monkey>();
        final var resultMap = new HashMap<String, Long>();
        expressionStack.push(new Monkey(monkeyName, rootEXpression));

        while (!expressionStack.isEmpty()) {
            final var monkey = expressionStack.pop();
            final var number = tryParse(monkey.expression());

            if (number == null) {
                SupportedOperations.getValues().stream()
                        .filter(delim -> monkey.expression().split("\\" + delim.value).length > 1).findFirst()
                        .ifPresent(sign -> {
                            final var exp = monkey.expression().split("\\" + sign.value);
                            final var left = exp[0].trim();
                            final var right = exp[1].trim();

                            final var leftNumber = resultMap.containsKey(left) ? resultMap.get(left) : tryParse(left);
                            final var rightNumber = resultMap.containsKey(right) ? resultMap.get(right)
                                    : tryParse(right);

                            if (leftNumber != null && rightNumber != null) {
                                final var result = eval(leftNumber, sign, rightNumber);
                                resultMap.put(monkey.name(), result);
                                return;
                            }

                            expressionStack.push(monkey);

                            if (leftNumber == null) {
                                expressionStack.push(new Monkey(left, monkeyMap.get(left)));
                            }

                            if (rightNumber == null) {
                                expressionStack.push(new Monkey(right, monkeyMap.get(right)));
                            }
                        });
            }

            if (number != null) {
                resultMap.put(monkey.name(), number);
            }
        }

        return resultMap.get(monkeyName);
    }

    private String eval(Map<String, String> monkeyMap) {
        final var rootEXpression = monkeyMap.get("root").split("\\+");
        final var leftEXpression = rootEXpression[0].trim();
        final var rightEXpression = rootEXpression[1].trim();

        monkeyMap.remove("humn");

        var rightSolution = eval(rightEXpression, monkeyMap.get(rightEXpression), monkeyMap);

        System.out.println(evalSingle(leftEXpression, monkeyMap));

        monkeyMap.put("humn", "3665520865940");

        System.out.println(
                "%s %s".formatted(eval(leftEXpression, monkeyMap.get(leftEXpression), monkeyMap).toString(),
                        eval(rightEXpression, monkeyMap.get(rightEXpression), monkeyMap).toString()));

        return monkeyMap.get("humn");
    }

    private SupportedOperations reverseSign(SupportedOperations sign) {
        return switch (sign) {
            case DIVIDE -> SupportedOperations.MULTIPLY;
            case MULTIPLY -> SupportedOperations.DIVIDE;
            case PLUS -> SupportedOperations.MINUS;
            case MINUS -> SupportedOperations.PLUS;
            default -> null;
        };
    }

    private String evalSingle(String monkeyName, Map<String, String> monkeyMap) {
        final var resultExpression = monkeyMap.get(monkeyName);
        var expression = "(%s)".formatted(resultExpression);
        final var resultSet = new LinkedList<>();

        final var expressionStack = new Stack<Monkey>();
        expressionStack.push(new Monkey(monkeyName, resultExpression));

        while (!expressionStack.isEmpty()) {

            final var monkey = expressionStack.pop();
            final var number = tryParse(monkey.expression());

            if (monkey.name() == "humn") {
                continue;
            }

            if (number == null) {
                final var sign = SupportedOperations.getValues().stream()
                        .filter(delim -> monkey.expression().split("\\" + delim.value).length > 1).findFirst()
                        .orElse(null);

                if (sign == null) {
                    continue;
                }

                final var exp = monkey.expression().split("\\" + sign.value);
                final var left = exp[0].trim();
                final var right = exp[1].trim();

                if (monkeyMap.containsKey(left)) {
                    final var isNum = tryParse(monkeyMap.get(left));
                    if (isNum == null) {
                        expression = expression.replace(left, "(%s)".formatted(monkeyMap.get(left)));
                    } else {
                        expression = expression.replace(left, monkeyMap.get(left));
                    }
                }

                if (monkeyMap.containsKey(right)) {
                    final var isNum = tryParse(monkeyMap.get(right));
                    if (isNum == null) {
                        expression = expression.replace(right, "(%s)".formatted(monkeyMap.get(right)));
                    } else {
                        expression = expression.replace(right, monkeyMap.get(right));
                    }

                }

                final var leftNumber = tryParse(left);
                final var rightNumber = tryParse(right);

                // if (leftNumber != null && rightNumber != null) {
                // // final var result = eval(leftNumber, sign, rightNumber);
                // continue;
                // }

                // expressionStack.push(monkey);

                if (leftNumber == null && monkeyMap.containsKey(left)) {
                    expressionStack.push(new Monkey(left, monkeyMap.get(left)));
                }

                if (rightNumber == null && monkeyMap.containsKey(right)) {
                    expressionStack.push(new Monkey(right, monkeyMap.get(right)));
                }
            }
        }

        return expression;
    }

    private Long findHumn(Map<String, Expression> reverseMap, Map<String, String> monkeyMap,
            Long otherSide) {
        // var expression = reverseMap.get(item);

        var items = new Stack<String>();
        items.push("humn");

        final var expressionStack = new LinkedList<Expression>();
        final var resultStack = monkeyMap.entrySet().stream().filter(s -> tryParse(s.getValue()) != null)
                .collect(Collectors.toMap(Map.Entry::getKey, s -> tryParse(s.getValue())));
        monkeyMap.remove("humn");

        final var passed = new HashSet<String>();
        while (!items.isEmpty()) {
            var item = items.pop();
            var expression = reverseMap.get(item);

            if (passed.contains(item)) {
                continue;
            }

            var itemNumber = tryParse(monkeyMap.get(item));
            if (itemNumber != null) {
                resultStack.put(item, itemNumber);
                continue;
            }

            // if (resultStack.containsKey(item) || passed.contains(item)) {
            // continue;
            // }

            // var itemNumber = tryParse(monkeyMap.get(item));
            // if (itemNumber != null) {
            // resultStack.put(item, itemNumber);
            // continue;
            // }

            var invertedExpression = invertExpression(expression, expression.left().equals(item));

            expressionStack.push(invertedExpression);

            items.push(invertedExpression.left());
            items.push(invertedExpression.right());

            passed.add(item);
            var left = tryParse(invertedExpression.left());
            var right = tryParse(invertedExpression.right());

            if (left == null) {
                left = tryParse(monkeyMap.get(invertedExpression.left()));
                if (left != null) {
                    resultStack.put(invertedExpression.left(), left);
                }
            }

            // if (right == null) {
            // right = tryParse(monkeyMap.get(invertedExpression.right()));
            // if (right != null) {
            // resultStack.put(invertedExpression.right(), right);
            // }
            // }
            // // // ? tryParse(invertedExpression.left())
            // // // : tryParse(monkeyMap.get(invertedExpression.left()));
            // // // ? tryParse(invertedExpression.right())
            // // // : tryParse(monkeyMap.get(invertedExpression.right()));

            // // invertedExpression = new Expression(invertedExpression.result(),
            // // left != null ? left.toString() : invertedExpression.left(),
            // invertedExpression.sign(),
            // // right != null ? right.toString() : invertedExpression.right());

            // if (left != null && right != null) {
            // resultStack.put(expression.result(), eval(left, expression.sign(), right));
            // } else {
            // expressionStack.push(invertedExpression);
            // }

            // if (left == null) {
            // items.push(invertedExpression.left());
            // }

            // if (right == null) {
            // items.push(invertedExpression.right());
            // }

            // passed.add(item);

            if (expression.result().equals("root")) {
                break;
            }
        }

        var rootExpression = expressionStack.pop();
        expressionStack.push(new Expression(rootExpression.result(), "0",
                reverseSign(rootExpression.sign()), otherSide.toString()));

        while (!expressionStack.isEmpty()) {
            var expression = expressionStack.poll();
            var left = tryParse(expression.left());
            // != null ? tryParse(expression.left())
            // : resultStack.get(expression.left());
            var right = tryParse(expression.right());
            // != null ? tryParse(expression.right())
            // : resultStack.get(expression.right());

            if (left == null) {
                left = tryParse(monkeyMap.get(expression.left()));
                if (left != null) {
                    resultStack.put(expression.left(), left);
                } else {
                    left = resultStack.get(expression.left());
                }
            }

            if (right == null) {
                right = tryParse(monkeyMap.get(expression.right()));
                if (right != null) {
                    resultStack.put(expression.right(), right);
                } else {
                    right = resultStack.get(expression.right());
                }
            }

            if (left != null && right != null) {
                resultStack.put(expression.result(), eval(left, expression.sign(), right));
            } else {
                expressionStack.offer(
                        new Expression(expression.result(), left != null ? left.toString() : expression.left(),
                                expression.sign(), right != null ? right.toString() : expression.right()));
            }

        }

        System.out.println(resultStack.get("humn"));
        System.out.println(expressionStack);

        return resultStack.get("humn");
    }

    private Expression invertExpression(Expression expression, boolean left) {
        return switch (expression.sign) {
            case PLUS -> left
                    ? new Expression(expression.left(), expression.result(), reverseSign(expression.sign()),
                            expression.right())
                    : new Expression(expression.right(), expression.result(), reverseSign(expression.sign()),
                            expression.left());

            case MINUS -> left
                    ? new Expression(expression.left(), expression.result(), reverseSign(expression.sign()),
                            expression.right())

                    : new Expression(expression.right(), expression.left(), expression.sign(),
                            expression.result());

            case MULTIPLY -> left
                    ? new Expression(expression.left(), expression.result(), reverseSign(expression.sign()),
                            expression.right())
                    : new Expression(expression.right(), expression.result(), reverseSign(expression.sign()),
                            expression.left());

            case DIVIDE -> left
                    ? new Expression(expression.left(), expression.result(), reverseSign(expression.sign()),
                            expression.right())
                    : new Expression(expression.right(), expression.left(), expression.sign(),
                            expression.result());

            default -> null;
        };

    }

    private Long tryParse(String expression) {
        try {
            return Long.parseLong(expression);
        } catch (Exception e) {
            return null;
        }
    }

    private Long eval(Long number1, SupportedOperations sign, Long number2) {
        return switch (sign) {
            case DIVIDE -> number1 / number2;
            case MULTIPLY -> number1 * number2;
            case PLUS -> number1 + number2;
            case MINUS -> number1 - number2;
            default -> null;
        };
    }

    public enum SupportedOperations {
        DIVIDE("/"),
        MULTIPLY("*"),
        PLUS("+"),
        MINUS("-");

        public final String value;

        private SupportedOperations(String value) {
            this.value = value;
        }

        public static List<SupportedOperations> getValues() {
            return List.of(SupportedOperations.values());
        }
    }

    public record Monkey(String name, String expression) {

    }

    public record Expression(String result, String left, SupportedOperations sign, String right) {

    }

    public record ExpressionEntry(String key, Expression expression) {

    }

    public static void main(String[] args) throws Exception {
        INSTANCE.run();
    }
}
