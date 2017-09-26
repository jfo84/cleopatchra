require_relative './data_mixin'

class Pull
  include Virtus.model
  include DataMixin

  attribute :repo_id, String

  def record
    connection.exec("INSERT INTO pulls (id, data, repo_id) VALUES ($1, $2, $3)", [id, data, repo_id])
  end

  def url
    data_hash['url']
  end
end