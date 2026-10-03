package ru.vorobevaqa.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;
import java.time.ZonedDateTime;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
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

  public String getUrl() {
    return "https://hub.docker.com/r/" + namespace + "/" + name;
  }
}
