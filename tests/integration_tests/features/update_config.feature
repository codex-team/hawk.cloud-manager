Feature: Update WireGuard config

  Scenario: Update config
    When I send PUT request to "/config" with data:
    """
    {
      "hosts": [
        {
          "name": "agent1",
          "public_key": "4Gm9s4KcYsSvikhJ0Tj4a60jQFI25QJVrYsJaQw0dgo=",
          "endpoint": "10.5.0.5:12345",
          "allowed_ips": ["10.5.0.0/32"]
        },
        {
          "name": "agent2",
          "public_key": "AZKH1M4ELjbgTLMgcf8rC4kf9CjS5qtXSn0xObgSqWc=",
          "endpoint": "10.5.0.6:12346",
          "allowed_ips": ["10.5.0.0/32"]
        }
      ],
      "groups": [
        {
          "name": "test-group",
          "hosts": ["agent1", "agent2"]
        }
      ]
    }
    """
    Then The response code should be 200
    When I send POST request to "/topology" with data:
    """
    {
      "public_key": "4Gm9s4KcYsSvikhJ0Tj4a60jQFI25QJVrYsJaQw0dgo="
    }
    """
    Then The response code should be 200
    And I receive data:
    """
    {
      "listen_port": 12345,
      "peers": [
        {
          "public_key": "AZKH1M4ELjbgTLMgcf8rC4kf9CjS5qtXSn0xObgSqWc=",
          "endpoint": "10.5.0.6:12346",
          "allowed_ips": ["10.5.0.0/32"]
        }
      ]
    }
    """
