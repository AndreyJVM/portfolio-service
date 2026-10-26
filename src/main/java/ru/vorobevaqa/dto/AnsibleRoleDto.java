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

@JsonIgnoreProperties(ignoreUnknown = true)
public class AnsibleRoleDto {

  private String name;
  private String description;

  @JsonProperty("download_count")
  private int downloadCount;

  @JsonProperty("github_user")
  private String githubUser;

  @JsonProperty("github_repo")
  private String githubRepo;

  public AnsibleRoleDto() {}

  public AnsibleRoleDto(
      String name, String description, int downloadCount, String githubUser, String githubRepo) {
    this.name = name;
    this.description = description;
    this.downloadCount = downloadCount;
    this.githubUser = githubUser;
    this.githubRepo = githubRepo;
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

  public int getDownloadCount() {
    return downloadCount;
  }

  public void setDownloadCount(int downloadCount) {
    this.downloadCount = downloadCount;
  }

  public String getGithubUser() {
    return githubUser;
  }

  public void setGithubUser(String githubUser) {
    this.githubUser = githubUser;
  }

  public String getGithubRepo() {
    return githubRepo;
  }

  public void setGithubRepo(String githubRepo) {
    this.githubRepo = githubRepo;
  }

  public String getUrl() {
    return "https://galaxy.ansible.com/ui/standalone/roles/" + githubUser + "/" + name + "/";
  }
}
