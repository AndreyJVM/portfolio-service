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

@JsonIgnoreProperties(ignoreUnknown = true)
public class DockerHubRepoDto {

  private String name;
  private String namespace;
  private String description;

  @JsonProperty("pull_count")
  private int pullCount;

  @JsonProperty("star_count")
  private int starCount;

  @JsonProperty("last_updated")
  private ZonedDateTime lastUpdated;

  public DockerHubRepoDto() {}

  public DockerHubRepoDto(
      String name,
      String namespace,
      String description,
      int pullCount,
      int starCount,
      ZonedDateTime lastUpdated) {
    this.name = name;
    this.namespace = namespace;
    this.description = description;
    this.pullCount = pullCount;
    this.starCount = starCount;
    this.lastUpdated = lastUpdated;
  }

  public String getName() {
    return name;
  }

  public void setName(String name) {
    this.name = name;
  }

  public String getNamespace() {
    return namespace;
  }

  public void setNamespace(String namespace) {
    this.namespace = namespace;
  }

  public String getDescription() {
    return description;
  }

  public void setDescription(String description) {
    this.description = description;
  }

  public int getPullCount() {
    return pullCount;
  }

  public void setPullCount(int pullCount) {
    this.pullCount = pullCount;
  }

  public int getStarCount() {
    return starCount;
  }

  public void setStarCount(int starCount) {
    this.starCount = starCount;
  }

  public ZonedDateTime getLastUpdated() {
    return lastUpdated;
  }

  public void setLastUpdated(ZonedDateTime lastUpdated) {
    this.lastUpdated = lastUpdated;
  }

  public String getUrl() {
    return "https://hub.docker.com/r/" + namespace + "/" + name;
  }

  public static Builder builder() {
    return new Builder();
  }

  public static class Builder {
    private String name;
    private String namespace;
    private String description;
    private int pullCount;
    private int starCount;
    private ZonedDateTime lastUpdated;

    public Builder name(String name) {
      this.name = name;
      return this;
    }

    public Builder namespace(String namespace) {
      this.namespace = namespace;
      return this;
    }

    public Builder description(String description) {
      this.description = description;
      return this;
    }

    public Builder pullCount(int pullCount) {
      this.pullCount = pullCount;
      return this;
    }

    public Builder starCount(int starCount) {
      this.starCount = starCount;
      return this;
    }

    public Builder lastUpdated(ZonedDateTime lastUpdated) {
      this.lastUpdated = lastUpdated;
      return this;
    }

    public DockerHubRepoDto build() {
      return new DockerHubRepoDto(name, namespace, description, pullCount, starCount, lastUpdated);
    }
  }
}
