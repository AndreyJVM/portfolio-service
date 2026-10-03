package ru.vorobevaqa.config;

import com.github.benmanes.caffeine.cache.Caffeine;
import java.time.Duration;
import java.util.List;
import org.springframework.cache.CacheManager;
import org.springframework.cache.caffeine.CaffeineCacheManager;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class CacheConfig {

  public static final String GITHUB_REPOS_CACHE = "githubRepos";
  public static final String DOCKER_REPOS_CACHE = "dockerRepos";

  @Bean
  public CacheManager cacheManager() {
    CaffeineCacheManager cacheManager = new CaffeineCacheManager();
    cacheManager.setCacheNames(List.of(GITHUB_REPOS_CACHE, DOCKER_REPOS_CACHE));
    cacheManager.setCaffeine(caffeineCacheBuilder());
    return cacheManager;
  }

  private Caffeine<Object, Object> caffeineCacheBuilder() {
    return Caffeine.newBuilder()
        .initialCapacity(10)
        .maximumSize(100)
        .expireAfterWrite(Duration.ofMinutes(60))
        .recordStats();
  }
}
