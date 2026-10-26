/*
 * Copyright 2024-2026 Andrey Vorobev
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */
package ru.vorobevaqa.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;
import java.time.ZonedDateTime;
import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
public class GitHubRepoDto {

  private String name;
  private String description;

  @JsonProperty("html_url")
  private String htmlUrl;

  private String language;
  private List<String> topics;

  @JsonProperty("stargazers_count")
  private int stars;

  @JsonProperty("forks_count")
  private int forks;

  @JsonProperty("updated_at")
  private ZonedDateTime updatedAt;

  private boolean fork;

  public GitHubRepoDto() {}

  public GitHubRepoDto(
      String name,
      String description,
      String htmlUrl,
      String language,
      List<String> topics,
      int stars,
      int forks,
      ZonedDateTime updatedAt,
      boolean fork) {
    this.name = name;
    this.description = description;
    this.htmlUrl = htmlUrl;
    this.language = language;
    this.topics = topics;
    this.stars = stars;
    this.forks = forks;
    this.updatedAt = updatedAt;
    this.fork = fork;
  }

  public String getName() {
    return name;
  }

  public void setName(String name) {
    this.name = name;
  }

  public String getDescription() {
    return description;
  }

  public void setDescription(String description) {
    this.description = description;
  }

  public String getHtmlUrl() {
    return htmlUrl;
  }

  public void setHtmlUrl(String htmlUrl) {
    this.htmlUrl = htmlUrl;
  }

  public String getLanguage() {
    return language;
  }

  public void setLanguage(String language) {
    this.language = language;
  }

  public List<String> getTopics() {
    return topics;
  }

  public void setTopics(List<String> topics) {
    this.topics = topics;
  }

  public int getStars() {
    return stars;
  }

  public void setStars(int stars) {
    this.stars = stars;
  }

  public int getForks() {
    return forks;
  }

  public void setForks(int forks) {
    this.forks = forks;
  }

  public ZonedDateTime getUpdatedAt() {
    return updatedAt;
  }

  public void setUpdatedAt(ZonedDateTime updatedAt) {
    this.updatedAt = updatedAt;
  }

  public boolean isFork() {
    return fork;
  }

  public void setFork(boolean fork) {
    this.fork = fork;
  }

  public static Builder builder() {
    return new Builder();
  }

  public static class Builder {
    private String name;
    private String description;
    private String htmlUrl;
    private String language;
    private List<String> topics;
    private int stars;
    private int forks;
    private ZonedDateTime updatedAt;
    private boolean fork;

    public Builder name(String name) {
      this.name = name;
      return this;
    }

    public Builder description(String description) {
      this.description = description;
      return this;
    }

    public Builder htmlUrl(String htmlUrl) {
      this.htmlUrl = htmlUrl;
      return this;
    }

    public Builder language(String language) {
      this.language = language;
      return this;
    }

    public Builder topics(List<String> topics) {
      this.topics = topics;
      return this;
    }

    public Builder stars(int stars) {
      this.stars = stars;
      return this;
    }

    public Builder forks(int forks) {
      this.forks = forks;
      return this;
    }

    public Builder updatedAt(ZonedDateTime updatedAt) {
      this.updatedAt = updatedAt;
      return this;
    }

    public Builder fork(boolean fork) {
      this.fork = fork;
      return this;
    }

    public GitHubRepoDto build() {
      return new GitHubRepoDto(
          name, description, htmlUrl, language, topics, stars, forks, updatedAt, fork);
    }
  }
}
