# --- Этап 1: Сборка артефакта ---
FROM maven:3.9-eclipse-temurin-21-alpine AS builder

WORKDIR /build

# Кэшируем зависимости (слой не пересобирается, пока не изменился pom.xml)
COPY pom.xml .
RUN mvn dependency:go-offline -B

# Копируем исходники и собираем jar без запуска тестов (тесты гоняются в CI)
COPY src ./src
RUN mvn clean package -DskipTests -B

# --- Этап 2: Production runtime ---
FROM eclipse-temurin:21-jre-alpine

RUN apk add --no-cache curl && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder --chown=appuser:appgroup /build/target/*.jar app.jar

USER appuser

EXPOSE 8080

ENV JAVA_OPTS="-XX:+UseContainerSupport -XX:MaxRAMPercentage=75.0"

ENTRYPOINT ["sh", "-c", "exec java $JAVA_OPTS -jar app.jar"]
