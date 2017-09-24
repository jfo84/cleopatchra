require_relative './data_mixin'

class Pull
  include Virtus.model
  include DataMixin

  attribute :repo_id, String

  def record
    connection.exec("INSERT INTO pulls (id, data, repo_id)
                    VALUES (#{id}, to_json('#{data}'::text), #{repo_id})")
  end

  def url
    data_hash['url']
  end
end